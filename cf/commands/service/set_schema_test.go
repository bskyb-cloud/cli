package service_test

import (
	"code.cloudfoundry.org/cli/cf/api/apifakes"
	"code.cloudfoundry.org/cli/cf/commandregistry"
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/models"
	"code.cloudfoundry.org/cli/cf/requirements"
	"code.cloudfoundry.org/cli/cf/requirements/requirementsfakes"
	testcmd "code.cloudfoundry.org/cli/utils/testhelpers/commands"
	testconfig "code.cloudfoundry.org/cli/utils/testhelpers/configuration"
	. "code.cloudfoundry.org/cli/utils/testhelpers/matchers"
	testterm "code.cloudfoundry.org/cli/utils/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("Set Schema for proxy service", func() {

	var args []string
	var requirementsFactory *requirementsfakes.FakeFactory

	var ui *testterm.FakeUI
	var serviceRepo *apifakes.FakeServiceRepository
	var config coreconfig.Repository
	var deps commandregistry.Dependency

	updateCommandDependency := func(pluginCall bool) {
		deps.UI = ui
		deps.Config = config
		deps.RepoLocator = deps.RepoLocator.SetServiceRepository(serviceRepo)
		commandregistry.Commands.SetCommand(commandregistry.Commands.FindCommand("set-schema").SetDependency(deps, pluginCall))
	}

	runCommand := func() bool {
		return testcmd.RunCLICommand("set-schema", args, requirementsFactory, updateCommandDependency, false, ui)
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		config = testconfig.NewRepositoryWithDefaults()
		serviceRepo = new(apifakes.FakeServiceRepository)
		requirementsFactory = new(requirementsfakes.FakeFactory)
	})

	Context("required params", func() {

		It("fails with usage when no args are passed", func() {
			args = []string{}
			requirementsFactory.NewLoginRequirementReturns(requirements.Passing{})
			requirementsFactory.NewTargetedSpaceRequirementReturns(requirements.Passing{})

			Expect(runCommand()).To(BeFalse())
		})

		It("fails with usage when only one arg is passed", func() {
			args = []string{"my-service"}
			requirementsFactory.NewLoginRequirementReturns(requirements.Passing{})
			requirementsFactory.NewTargetedSpaceRequirementReturns(requirements.Passing{})

			Expect(runCommand()).To(BeFalse())
		})

		It("fails when not logged in", func() {
			args = []string{"my-service", "schema-file"}
			requirementsFactory.NewLoginRequirementReturns(requirements.Failing{})
			requirementsFactory.NewTargetedSpaceRequirementReturns(requirements.Failing{})

			Expect(runCommand()).To(BeFalse())
		})

		It("fails when logged in but space is not targetted", func() {
			args = []string{"my-service", "schema-file"}
			requirementsFactory.NewLoginRequirementReturns(requirements.Passing{})
			requirementsFactory.NewTargetedSpaceRequirementReturns(requirements.Failing{})

			Expect(runCommand()).To(BeFalse())
		})

		It("succeeds when logged in with targetted space and two arguments are passed", func() {
			args = []string{"my-service", "schemafile.txt"}

			serviceInstance := models.ServiceInstance{}
			serviceInstance.Name = "my-service"
			serviceReq := new(requirementsfakes.FakeServiceInstanceRequirement)
			serviceReq.GetServiceInstanceReturns(serviceInstance)

			requirementsFactory.NewLoginRequirementReturns(requirements.Passing{})
			requirementsFactory.NewTargetedSpaceRequirementReturns(requirements.Passing{})
			requirementsFactory.NewServiceInstanceRequirementReturns(serviceReq)

			ioutil.WriteFile("schemafile.txt", []byte(nil), 0644)

			Expect(runCommand()).To(BeTrue())

			os.Remove("schemafile.txt")
		})

	})

	Context("schema file", func() {

		BeforeEach(func() {
			serviceInstance := models.ServiceInstance{}
			serviceInstance.Name = "my-proxy-service"
			serviceReq := new(requirementsfakes.FakeServiceInstanceRequirement)
			serviceReq.GetServiceInstanceReturns(serviceInstance)

			requirementsFactory.NewLoginRequirementReturns(requirements.Passing{})
			requirementsFactory.NewTargetedSpaceRequirementReturns(requirements.Passing{})
			requirementsFactory.NewServiceInstanceRequirementReturns(serviceReq)

			args = []string{"my-proxy-service", "schemafile.txt"}
		})

		AfterEach(func() {
			os.Remove("schemafile.txt")
		})

		It("fails when schema file does not exists", func() {
			Expect(runCommand()).To(BeFalse())

			Expect(ui.Outputs()).To(ContainSubstrings(
				[]string{"Failed to read file schemafile.txt."},
				[]string{"Error: open schemafile.txt: no such file or directory"},
			))
		})

		It("works when empty schema file is given", func() {
			ioutil.WriteFile("schemafile.txt", []byte(nil), 0644)
			Expect(runCommand()).To(BeTrue())

			Expect(ui.Outputs()).To(ContainSubstrings(
				[]string{"Applying schema to my-proxy-service in org my-org / space my-space as my-user..."},
				[]string{"OK"},
			))
		})

		It("works when schema file with no duplicates is given", func() {
			ioutil.WriteFile("schemafile.txt", []byte(".sky.com"), 0644)
			Expect(runCommand()).To(BeTrue())

			Expect(ui.Outputs()).To(ContainSubstrings(
				[]string{"Applying schema to my-proxy-service in org my-org / space my-space as my-user..."},
				[]string{"OK"},
			))
		})

		It("fails when schema file with duplicates is given", func() {
			ioutil.WriteFile("schemafile.txt", []byte(".sky.com\nupload.sky.com"), 0644)
			Expect(runCommand()).To(BeFalse())

			Expect(ui.Outputs()).To(ContainSubstrings(
				[]string{`Failed: top level domain 'sky.com' duplicates found: [".sky.com" "upload.sky.com"]`},
			))
		})

	})

})
