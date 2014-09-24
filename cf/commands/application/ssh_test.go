package application_test

import (
	. "github.com/cloudfoundry/cli/cf/commands/application"
	"github.com/cloudfoundry/cli/cf/models"
	testapi "github.com/cloudfoundry/cli/testhelpers/api"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing with ginkgo", func() {

	// stub out the actual ssh command
	ExecuteCmd = func(appname string, args []string) (err error) {
		return
	}

	It("TestSshRequirements", func() {
		args := []string{"my-app"}

		ExecuteCmd = func(appname string, args []string) (err error) {
			return
		}

		appSshRepo := &testapi.FakeAppSshRepo{}

		reqFactory := &testreq.FakeReqFactory{LoginSuccess: false, TargetedSpaceSuccess: true, Application: models.Application{}}
		callSsh(args, reqFactory, appSshRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())

		reqFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: false, Application: models.Application{}}
		callSsh(args, reqFactory, appSshRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())

		reqFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: true, Application: models.Application{}}
		callSsh(args, reqFactory, appSshRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeTrue())
		Expect(reqFactory.ApplicationName).To(Equal("my-app"))

	})

	It("TestSshFailsWithUsage", func() {

		appFilesRepo := &testapi.FakeAppSshRepo{}
		reqFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: true, Application: models.Application{}}
		ui := callSsh([]string{}, reqFactory, appFilesRepo)

		Expect(ui.FailedWithUsage).To(BeTrue())
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
	})

	It("TestGettingSshDetails", func() {

		app := models.Application{}
		app.Name = "my-found-app"
		app.Guid = "my-app-guid"

		reqFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: true, Application: app}

		var sshInfo models.SshConnectionDetails
		sshInfo.Ip = "10.0.0.1"
		sshInfo.Port = 1234
		sshInfo.User = "vcap"
		sshInfo.SshKey = "fakekey"

		appSshRepo := &testapi.FakeAppSshRepo{SshDetails: sshInfo}

		ui := callSsh([]string{"my-app"}, reqFactory, appSshRepo)

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"SSHing to application my-found-app, instance 0..."},
			[]string{"OK"},
			[]string{"SSH username is vcap"},
			[]string{"SSH IP Address is 10.0.0.1"},
			[]string{"SSH Port is 1234"},
			[]string{"SSH Identity is"},
			[]string{"Command:"},
			[]string{"SSH Finished"},
		))

		Expect(appSshRepo.AppGuid).To(Equal("my-app-guid"))
	})
})

func callSsh(args []string, reqFactory *testreq.FakeReqFactory, appSshRepo *testapi.FakeAppSshRepo) (ui *testterm.FakeUI) {

	ui = &testterm.FakeUI{}

	configRepo := testconfig.NewRepositoryWithDefaults()
	cmd := NewSsh(ui, configRepo, appSshRepo)
	testcmd.RunCommand(cmd, args, reqFactory)

	return
}
