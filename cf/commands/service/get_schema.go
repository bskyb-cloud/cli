package service

import (
	"github.com/cloudfoundry/cli/cf"
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/command_metadata"
	"github.com/cloudfoundry/cli/cf/configuration"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/codegangsta/cli"
)

type GetSchema struct {
	ui                 terminal.UI
	config             configuration.Reader
	serviceRepo        api.ServiceRepository
	serviceInstanceReq requirements.ServiceInstanceRequirement
}

func NewGetSchema(ui terminal.UI, config configuration.Reader, serviceRepo api.ServiceRepository) (cmd *GetSchema) {
	cmd = new(GetSchema)
	cmd.ui = ui
	cmd.config = config
	cmd.serviceRepo = serviceRepo
	return
}

func (command *GetSchema) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "get-schema",
		ShortName:   "gs",
		Description: "Get a service schema. Currently only supported in the webproxy.",
		Usage:       "CF_NAME get-schema SERVICE_INSTANCE",
	}
}

func (cmd *GetSchema) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 1 {
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

func (cmd *GetSchema) Run(c *cli.Context) {
	serviceInstance := cmd.serviceInstanceReq.GetServiceInstance()

	cmd.ui.Say("Getting schema for %s in org %s / space %s as %s...",
		terminal.EntityNameColor(serviceInstance.Name),
		terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
		terminal.EntityNameColor(cmd.config.SpaceFields().Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)
	schema, err := cmd.serviceRepo.GetSchema(serviceInstance)

	if err != nil {
		if httpError, ok := err.(errors.HttpError); ok && httpError.ErrorCode() == errors.SERVICE_INSTANCE_NAME_TAKEN {
			cmd.ui.Failed("%s\nTIP: Use '%s services' to view all services in this org and space.", httpError.Error(), cf.Name())
		} else {
			cmd.ui.Failed(err.Error())
		}
	}

	cmd.ui.Ok()
	cmd.ui.Say("Schema is:\n\n%s\n", schema)
}
