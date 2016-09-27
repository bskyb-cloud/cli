package service

import (
	"code.cloudfoundry.org/cli/cf"
	"code.cloudfoundry.org/cli/cf/api"
	"code.cloudfoundry.org/cli/cf/commandregistry"
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/cf/flags"
	"code.cloudfoundry.org/cli/cf/requirements"
	"code.cloudfoundry.org/cli/cf/terminal"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"strings"
)

type SetSchema struct {
	ui                 terminal.UI
	config             coreconfig.Reader
	serviceRepo        api.ServiceRepository
	serviceInstanceReq requirements.ServiceInstanceRequirement
}

func init() {
	commandregistry.Register(&SetSchema{})
}

func (cmd *SetSchema) MetaData() commandregistry.CommandMetadata {
	return commandregistry.CommandMetadata{
		Name:        "set-schema",
		ShortName:   "ss",
		Description: "Set schema for a service. Currently only supported in the webproxy.",
		Usage:       []string{"CF_NAME set-schema SERVICE_INSTANCE SCHEME_FILENAME"},
	}
}

func (cmd *SetSchema) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) (reqs []requirements.Requirement, err error) {

	if len(fc.Args()) != 2 {
		cmd.ui.Failed("Incorrect Usage." + "\n\n" + commandregistry.Commands.CommandUsage("set-schema"))
		return nil, fmt.Errorf("Incorrect usage: %d arguments of %d required", len(fc.Args()), 2)
	}

	cmd.serviceInstanceReq = requirementsFactory.NewServiceInstanceRequirement(fc.Args()[0])

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedSpaceRequirement(),
		cmd.serviceInstanceReq,
	}

	return reqs, nil
}

func (cmd *SetSchema) SetDependency(deps commandregistry.Dependency, pluginCall bool) commandregistry.Command {
	cmd.ui = deps.UI
	cmd.config = deps.Config
	cmd.serviceRepo = deps.RepoLocator.GetServiceRepository()
	return cmd
}

func (cmd *SetSchema) Execute(fc flags.FlagContext) error {
	schemaFilename := fc.Args()[1]

	schemaBytes, ferr := ioutil.ReadFile(schemaFilename)
	if ferr != nil {
		cmd.ui.Failed("Failed to read file %s. Error: %s", schemaFilename, ferr)
		return errors.New("Error reading schema file: " + ferr.Error())
	}

	schema := string(schemaBytes)

	err := validateForDuplicates(schema)
	if err != nil {
		cmd.ui.Failed(err.Error())
		return errors.New("Schema file invalid: " + err.Error())
	}

	serviceInstance := cmd.serviceInstanceReq.GetServiceInstance()

	cmd.ui.Say("Applying schema to %s in org %s / space %s as %s...",
		terminal.EntityNameColor(serviceInstance.Name),
		terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
		terminal.EntityNameColor(cmd.config.SpaceFields().Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	err = cmd.serviceRepo.SetSchema(serviceInstance, schema)

	if err != nil {
		if httpError, ok := err.(errors.HTTPError); ok && httpError.ErrorCode() == errors.ServiceInstanceNameTaken {
			cmd.ui.Failed("%s\nTIP: Use '%s services' to view all services in this org and space.", httpError.Error(), cf.Name)
		} else {
			cmd.ui.Failed(err.Error())
		}
		return errors.New("Error setting schema: " + err.Error())
	}

	cmd.ui.Ok()
	return nil
}

func validateForDuplicates(fileContents string) error {
	lines := strings.Split(fileContents, "\n")
	tlds := make(map[string][]string)

	for _, line := range lines {
		if line == "" {
			continue
		}

		tld, e := publicsuffix.EffectiveTLDPlusOne(strings.TrimSuffix(line, "/"))
		if e != nil {
			return e
		}

		tlds[tld] = append(tlds[tld], line)

		if len(tlds[tld]) > 1 {
			return errors.New(fmt.Sprintf("Failed: top level domain '%v' duplicates found: %q", tld, tlds[tld]))
		}
	}
	return nil
}
