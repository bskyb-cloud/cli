package service

import (
	"fmt"
	"github.com/cloudfoundry/cli/cf"
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/command_metadata"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/codegangsta/cli"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"strings"
)

type SetSchema struct {
	ui                 terminal.UI
	config             core_config.Reader
	serviceRepo        api.ServiceRepository
	serviceInstanceReq requirements.ServiceInstanceRequirement
}

func NewSetSchema(ui terminal.UI, config core_config.Reader, serviceRepo api.ServiceRepository) (cmd *SetSchema) {
	cmd = new(SetSchema)
	cmd.ui = ui
	cmd.config = config
	cmd.serviceRepo = serviceRepo
	return
}

func (command *SetSchema) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "set-schema",
		ShortName:   "ss",
		Description: "Set schema for a service. Currently only supported in the webproxy.",
		Usage:       "CF_NAME set-schema SERVICE_INSTANCE SCHEME_FILENAME",
	}
}

func (cmd *SetSchema) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 2 {
		err = errors.New("incorrect usage")
		cmd.ui.FailWithUsage(c)
		return
	}

	cmd.serviceInstanceReq = requirementsFactory.NewServiceInstanceRequirement(c.Args()[0])

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedSpaceRequirement(),
		cmd.serviceInstanceReq,
	}

	return
}

func (cmd *SetSchema) Run(c *cli.Context) {
	schemaFilename := c.Args()[1]

	schemaBytes, ferr := ioutil.ReadFile(schemaFilename)
	if ferr != nil {
		cmd.ui.Failed("Failed to read file %s. Error: %s", schemaFilename, ferr)
		return
	}

	schema := string(schemaBytes)

	err := validateForDuplicates(schema)
	if err != nil {
		cmd.ui.Failed(err.Error())
		return
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
		if httpError, ok := err.(errors.HttpError); ok && httpError.ErrorCode() == errors.SERVICE_INSTANCE_NAME_TAKEN {
			cmd.ui.Failed("%s\nTIP: Use '%s services' to view all services in this org and space.", httpError.Error(), cf.Name())
		} else {
			cmd.ui.Failed(err.Error())
		}
	}

	cmd.ui.Ok()
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
