package application

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"

	"code.cloudfoundry.org/cli/cf/api"
	"code.cloudfoundry.org/cli/cf/commandregistry"
	"code.cloudfoundry.org/cli/cf/commands"
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/flags"
	. "code.cloudfoundry.org/cli/cf/i18n"
	"code.cloudfoundry.org/cli/cf/models"
	"code.cloudfoundry.org/cli/cf/net"
	"code.cloudfoundry.org/cli/cf/requirements"
	sshCmd "code.cloudfoundry.org/cli/cf/ssh"
	"code.cloudfoundry.org/cli/cf/ssh/options"
	sshTerminal "code.cloudfoundry.org/cli/cf/ssh/terminal"
	"code.cloudfoundry.org/cli/cf/terminal"
	"errors"
	"io/ioutil"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type SSH struct {
	ui            terminal.UI
	config        coreconfig.Reader
	gateway       net.Gateway
	appReq        requirements.ApplicationRequirement
	sshCodeGetter commands.SSHCodeGetter
	opts          *options.SSHOptions
	secureShell   sshCmd.SecureShell

	appSshRepo api.AppSshRepository // nimbus DEA
}

type sshInfo struct {
	SSHEndpoint            string `json:"app_ssh_endpoint"`
	SSHEndpointFingerprint string `json:"app_ssh_host_key_fingerprint"`
}

func init() {
	commandregistry.Register(&SSH{})
}

func (cmd *SSH) MetaData() commandregistry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["L"] = &flags.StringSliceFlag{ShortName: "L", Usage: T("Local port forward specification. This flag can be defined more than once.")}
	fs["command"] = &flags.StringSliceFlag{Name: "command", ShortName: "c", Usage: T("Command to run. This flag can be defined more than once.")}
	fs["app-instance-index"] = &flags.IntFlag{Name: "app-instance-index", ShortName: "i", Usage: T("Application instance index")}
	fs["skip-host-validation"] = &flags.BoolFlag{Name: "skip-host-validation", ShortName: "k", Usage: T("Skip host key validation")}
	fs["skip-remote-execution"] = &flags.BoolFlag{Name: "skip-remote-execution", ShortName: "N", Usage: T("Do not execute a remote command")}
	fs["request-pseudo-tty"] = &flags.BoolFlag{Name: "request-pseudo-tty", ShortName: "t", Usage: T("Request pseudo-tty allocation")}
	fs["force-pseudo-tty"] = &flags.BoolFlag{Name: "force-pseudo-tty", ShortName: "tt", Usage: T("Force pseudo-tty allocation")}
	fs["disable-pseudo-tty"] = &flags.BoolFlag{Name: "disable-pseudo-tty", ShortName: "T", Usage: T("Disable pseudo-tty allocation")}

	return commandregistry.CommandMetadata{
		Name:        "ssh",
		Description: T("SSH to an application container instance"),
		Usage: []string{
			T("CF_NAME ssh APP_NAME [-i app-instance-index] [-c command] [-L [bind_address:]port:host:hostport] [--skip-host-validation] [--skip-remote-execution] [--request-pseudo-tty] [--force-pseudo-tty] [--disable-pseudo-tty]"),
		},
		Flags: fs,
	}
}

func (cmd *SSH) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) ([]requirements.Requirement, error) {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(T("Incorrect Usage. Requires APP_NAME as argument") + "\n\n" + commandregistry.Commands.CommandUsage("ssh"))
		return nil, fmt.Errorf("Incorrect usage: %d arguments of %d required", len(fc.Args()), 1)
	}

	if fc.IsSet("i") && fc.Int("i") < 0 {
		cmd.ui.Failed(fmt.Sprintf(T("Incorrect Usage:")+" %s\n\n%s", T("Value for flag 'app-instance-index' cannot be negative"), commandregistry.Commands.CommandUsage("ssh")))
		return nil, fmt.Errorf("Incorrect usage: app-instance-index cannot be negative")
	}

	var err error
	cmd.opts, err = options.NewSSHOptions(fc)

	if err != nil {
		cmd.ui.Failed(fmt.Sprintf(T("Incorrect Usage:")+" %s\n\n%s", err.Error(), commandregistry.Commands.CommandUsage("ssh")))
		return nil, err
	}

	cmd.appReq = requirementsFactory.NewApplicationRequirement(cmd.opts.AppName)

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedSpaceRequirement(),
		cmd.appReq,
	}

	return reqs, nil
}

func (cmd *SSH) SetDependency(deps commandregistry.Dependency, pluginCall bool) commandregistry.Command {
	cmd.ui = deps.UI
	cmd.config = deps.Config
	cmd.gateway = deps.Gateways["cloud-controller"]

	if deps.WildcardDependency != nil {
		cmd.secureShell = deps.WildcardDependency.(sshCmd.SecureShell)
	}

	//get ssh-code for dependency
	sshCodeGetter := commandregistry.Commands.FindCommand("ssh-code")
	sshCodeGetter = sshCodeGetter.SetDependency(deps, false)
	cmd.sshCodeGetter = sshCodeGetter.(commands.SSHCodeGetter)

	cmd.appSshRepo = deps.RepoLocator.GetAppSshRepository() // nimbus DEA

	return cmd
}

func (cmd *SSH) Execute(fc flags.FlagContext) error {
	app := cmd.appReq.GetApplication()

	if app.Diego {
		return cmd.executeDiego(fc, app)
	} else {
		return cmd.executeNimbusDEA(fc, app)
	}
}

func (cmd *SSH) executeDiego(fc flags.FlagContext, app models.Application) error {
	//app := cmd.appReq.GetApplication()
	info, err := cmd.getSSHEndpointInfo()
	if err != nil {
		return errors.New(T("Error getting SSH info:") + err.Error())
	}

	sshAuthCode, err := cmd.sshCodeGetter.Get()
	if err != nil {
		return errors.New(T("Error getting one time auth code: ") + err.Error())
	}

	//init secureShell if it is not already set by SetDependency() with fakes
	if cmd.secureShell == nil {
		cmd.secureShell = sshCmd.NewSecureShell(
			sshCmd.DefaultSecureDialer(),
			sshTerminal.DefaultHelper(),
			sshCmd.DefaultListenerFactory(),
			30*time.Second,
			app,
			info.SSHEndpointFingerprint,
			info.SSHEndpoint,
			sshAuthCode,
		)
	}

	err = cmd.secureShell.Connect(cmd.opts)
	if err != nil {
		return errors.New(T("Error opening SSH connection: ") + err.Error())
	}
	defer cmd.secureShell.Close()

	err = cmd.secureShell.LocalPortForward()
	if err != nil {
		return errors.New(T("Error forwarding port: ") + err.Error())
	}

	if cmd.opts.SkipRemoteExecution {
		err = cmd.secureShell.Wait()
	} else {
		err = cmd.secureShell.InteractiveSession()
	}

	if err != nil {
		if exitError, ok := err.(*ssh.ExitError); ok {
			exitStatus := exitError.ExitStatus()
			if sig := exitError.Signal(); sig != "" {
				cmd.ui.Say(T("Process terminated by signal: {{.Signal}}. Exited with {{.ExitCode}}", map[string]interface{}{
					"Signal":   sig,
					"ExitCode": exitStatus,
				}))
			}
			os.Exit(exitStatus)
		} else {
			return errors.New(T("Error: ") + err.Error())
		}
	}
	return nil
}

func (cmd *SSH) getSSHEndpointInfo() (sshInfo, error) {
	info := sshInfo{}
	err := cmd.gateway.GetResource(cmd.config.APIEndpoint()+"/v2/info", &info)
	return info, err
}

var ExecuteCmd = func(appname string, args []string) (err error) {

	command := exec.Command(appname, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	return
}

func (cmd *SSH) executeNimbusDEA(fc flags.FlagContext, app models.Application) error {

	instance := int(cmd.opts.Index)
	sshapi := cmd.appSshRepo

	cmd.ui.Say("SSHing to application %s, instance %s...",
		terminal.EntityNameColor(app.Name),
		terminal.EntityNameColor(strconv.Itoa(instance)),
	)

	sshDetails, apiErr := sshapi.GetSshDetails(app.GUID, instance)

	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return errors.New("Error getting ssh details: " + apiErr.Error())
	}

	cmd.ui.Ok()

	tempdir, error := ioutil.TempDir("", "gocf")
	if error != nil {
		return errors.New("Error tmp dir: " + error.Error())
	}

	tempfile := path.Join(tempdir, "identity")

	error = ioutil.WriteFile(tempfile, []byte(sshDetails.SshKey), 0600)
	if error != nil {
		return errors.New("Error writing ssh keys to tmp file: " + error.Error())
	}

	cmd.ui.Say("SSH username is %s", terminal.EntityNameColor(sshDetails.User))
	cmd.ui.Say("SSH IP Address is %s", terminal.EntityNameColor(sshDetails.Ip))
	cmd.ui.Say("SSH Port is %s", terminal.EntityNameColor(strconv.Itoa(sshDetails.Port)))
	cmd.ui.Say("SSH Identity is %s", terminal.EntityNameColor(tempfile))

	cmd.ui.Say("")

	userAndHost := fmt.Sprintf("%s@%s", sshDetails.User, sshDetails.Ip)
	var sshCommand []string
	sshCommand = []string{"-i", tempfile, "-o", "ConnectTimeout=5", "-o", "StrictHostKeychecking=no", "-o", "UserKnownHostsFile=/dev/null", "-p", strconv.Itoa(sshDetails.Port), userAndHost}

	cmd.ui.Say("Command: ssh %s", strings.Join(sshCommand, " "))

	err := ExecuteCmd("ssh", sshCommand)
	if err != nil {
		cmd.ui.Say("Command Failed: %s", err)
		return errors.New("Error running ssh command: " + err.Error())
	}

	err2 := os.Remove(tempfile)
	if err2 != nil {
		panic(err2)
	}

	err3 := os.Remove(tempdir)
	if err3 != nil {
		panic(err3)
	}

	cmd.ui.Say("SSH Finished\n")
	return nil
}
