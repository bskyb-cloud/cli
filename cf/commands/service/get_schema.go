package service

import (
	"github.com/cloudfoundry/cli/cf"
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/commandregistry"
	"github.com/cloudfoundry/cli/cf/configuration/coreconfig"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/cloudfoundry/cli/cf/flags"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
)

type GetSchema struct {
	ui                 terminal.UI
	config             coreconfig.Reader
	serviceRepo        api.ServiceRepository
	serviceInstanceReq requirements.ServiceInstanceRequirement
}

func init() {
	commandregistry.Register(&GetSchema{})
}

func (cmd *GetSchema) MetaData() commandregistry.CommandMetadata {
	return commandregistry.CommandMetadata{
		Name:        "get-schema",
		ShortName:   "gs",
		Description: "Get a service schema. Currently only supported in the webproxy.",
		Usage:       "CF_NAME get-schema SERVICE_INSTANCE",
	}
}

func (cmd *GetSchema) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) (reqs []requirements.Requirement) {

	if len(fc.Args()) != 1 {
		cmd.ui.Failed("Incorrect Usage." + "\n\n" + commandregistry.Commands.CommandUsage("get-schema"))
	}

	cmd.serviceInstanceReq = requirementsFactory.NewServiceInstanceRequirement(fc.Args()[0])

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedSpaceRequirement(),
		cmd.serviceInstanceReq,
	}

	return
}

func (cmd *GetSchema) SetDependency(deps commandregistry.Dependency, pluginCall bool) commandregistry.Command {
	cmd.ui = deps.UI
	cmd.config = deps.Config
	cmd.serviceRepo = deps.RepoLocator.GetServiceRepository()
	return cmd
}

func (cmd *GetSchema) Execute(fc flags.FlagContext) error {
	serviceInstance := cmd.serviceInstanceReq.GetServiceInstance()

	cmd.ui.Say("Getting schema for %s in org %s / space %s as %s...",
		terminal.EntityNameColor(serviceInstance.Name),
		terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
		terminal.EntityNameColor(cmd.config.SpaceFields().Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)
	schema, err := cmd.serviceRepo.GetSchema(serviceInstance)

	if err != nil {
		if httpError, ok := err.(errors.HTTPError); ok && httpError.ErrorCode() == errors.SERVICE_INSTANCE_NAME_TAKEN {
			cmd.ui.Failed("%s\nTIP: Use '%s services' to view all services in this org and space.", httpError.Error(), cf.Name())
		} else {
			cmd.ui.Failed(err.Error())
		}
		return errors.New("Error reading schema: " + err.Error())
	}

	cmd.ui.Ok()
	cmd.ui.Say("Schema is:\n\n%s\n", schema)
	return nil
}
