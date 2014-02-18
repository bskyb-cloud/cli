package application

import (
	"cf/api"
	"cf/configuration"
	"cf/requirements"
	"cf/terminal"
	"errors"
	"strconv"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"path"
	"fmt"
	"os/exec"
    "os"
    "strings"
)

type Ssh struct {
	ui           terminal.UI
	config       configuration.Reader
	appSshRepo   api.AppSshRepository
	appReq       requirements.ApplicationRequirement
}

func NewSsh(ui terminal.UI, config configuration.Reader, appSshRepo api.AppSshRepository) (cmd *Ssh) {
	cmd = new(Ssh)
	cmd.ui = ui
	cmd.config = config
	cmd.appSshRepo = appSshRepo
	return
}

func (cmd *Ssh) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) < 1 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c, "ssh")
		return
	}

	cmd.appReq = reqFactory.NewApplicationRequirement(c.Args()[0])

	reqs = []requirements.Requirement{
		reqFactory.NewLoginRequirement(),
		reqFactory.NewTargetedSpaceRequirement(),
		cmd.appReq,
	}
	return
}

func (cmd *Ssh) Run(c *cli.Context) {
	app := cmd.appReq.GetApplication()

	instance := c.Int("instance")
	
	sshapi:= cmd.appSshRepo

	cmd.ui.Say("SSHing to application %s, instance %s...",
		terminal.EntityNameColor(app.Name), 
		terminal.EntityNameColor(strconv.Itoa(instance)),
	)

	
    apiResponse, sshDetails := sshapi.GetSshDetails(app.Guid, instance)
	if apiResponse.IsNotSuccessful() {
		cmd.ui.Failed(apiResponse.Message)
		return
	}
	
	cmd.ui.Ok()
	
	tempdir, error := ioutil.TempDir("", "gocf")
	if error != nil { panic(error) }
	
	tempfile := path.Join(tempdir, "identity")
	
	error = ioutil.WriteFile(tempfile, []byte(sshDetails.SshKey), 0600)
	if error != nil { panic(error) }
	
	cmd.ui.Say("SSH username is %s", terminal.EntityNameColor(sshDetails.User))
	cmd.ui.Say("SSH IP Address is %s", terminal.EntityNameColor(sshDetails.Ip))
    cmd.ui.Say("SSH Port is %s", terminal.EntityNameColor(strconv.Itoa(sshDetails.Port)))
    cmd.ui.Say("SSH Identity is %s", terminal.EntityNameColor(tempfile))
	
	cmd.ui.Say("")
	
	userAndHost := fmt.Sprintf("%s@%s", sshDetails.User, sshDetails.Ip)
	var sshCommand []string
	sshCommand = []string{"-i", tempfile, "-o", "ConnectTimeout=5", "-o", "StrictHostKeychecking=no", "-o", "UserKnownHostsFile=/dev/null", "-p", strconv.Itoa(sshDetails.Port), userAndHost}
	
	cmd.ui.Say("Command: %s", strings.Join(sshCommand," "))

	// Check ssh is in path
	
	command := exec.Command("ssh", sshCommand...)
	
	command.Stdin = os.Stdin 
    command.Stdout = os.Stdout 
    command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
	  fmt.Printf("%s\n", err) 
	  panic(err) 
	} 
	
	err2 := os.Remove(tempfile)
	if err2 != nil { panic(err2) }
	
	err3 := os.Remove(tempdir)
	if err3 != nil { panic(err3) }
	
	cmd.ui.Say("SSH Finished\n")
}
