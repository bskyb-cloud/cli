package service

import (
	"github.com/nimbus-cloud/cli/cf"
	"github.com/nimbus-cloud/cli/cf/api"
	"github.com/nimbus-cloud/cli/cf/command_metadata"
	"github.com/nimbus-cloud/cli/cf/configuration"
	"github.com/nimbus-cloud/cli/cf/errors"
	"github.com/nimbus-cloud/cli/cf/requirements"
	"github.com/nimbus-cloud/cli/cf/terminal"
	"github.com/codegangsta/cli"
)

type ApplySchema struct {
	ui                 terminal.UI
	config             configuration.Reader
	serviceRepo        api.ServiceRepository
	serviceInstanceReq requirements.ServiceInstanceRequirement
}

func NewApplySchema(ui terminal.UI, config configuration.Reader, serviceRepo api.ServiceRepository) (cmd *ApplySchema) {
	cmd = new(ApplySchema)
	cmd.ui = ui
	cmd.config = config
	cmd.serviceRepo = serviceRepo
	return
}

func (command *ApplySchema) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "apply-schema",
		ShortName:   "as",
		Description: "Apply a schema to a service. Currently only supported in the webproxy.",
		Usage:       "CF_NAME apply-schema SERVICE_INSTANCE SCHEME",
	}
}

func (cmd *ApplySchema) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
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

func (cmd *ApplySchema) Run(c *cli.Context) {
	schema := c.Args()[1]
	serviceInstance := cmd.serviceInstanceReq.GetServiceInstance()

	cmd.ui.Say("Applying schema to %s in org %s / space %s as %s...",
		terminal.EntityNameColor(serviceInstance.Name),
		terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
		terminal.EntityNameColor(cmd.config.SpaceFields().Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)
	err := cmd.serviceRepo.ApplySchema(serviceInstance, schema)

	if err != nil {
		if httpError, ok := err.(errors.HttpError); ok && httpError.ErrorCode() == errors.SERVICE_INSTANCE_NAME_TAKEN {
			cmd.ui.Failed("%s\nTIP: Use '%s services' to view all services in this org and space.", httpError.Error(), cf.Name())
		} else {
			cmd.ui.Failed(err.Error())
		}
	}

	cmd.ui.Ok()
}
