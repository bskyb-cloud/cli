package commands

import (
	"github.com/nimbus-cloud/cli/cf/models"
)

type FakeAppDisplayer struct {
	AppToDisplay models.Application
}

func (displayer *FakeAppDisplayer) ShowApp(app models.Application) {
	displayer.AppToDisplay = app
}
