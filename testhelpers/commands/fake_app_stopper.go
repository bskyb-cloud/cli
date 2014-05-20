package commands

import (
	"github.com/nimbus-cloud/cli/cf/models"
)

type FakeAppStopper struct {
	AppToStop models.Application
}

func (stopper *FakeAppStopper) ApplicationStop(app models.Application) (updatedApp models.Application, err error) {
	stopper.AppToStop = app
	updatedApp = app
	return
}
