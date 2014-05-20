package commands

import (
	"github.com/nimbus-cloud/cli/cf/models"
)

type FakeAppRestarter struct {
	AppToRestart models.Application
}

func (restarter *FakeAppRestarter) ApplicationRestart(appToRestart models.Application) {
	restarter.AppToRestart = appToRestart
	return
}
