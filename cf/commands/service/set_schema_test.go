package service_test

import (
	testapi "github.com/cloudfoundry/cli/cf/api/fakes"
	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/models"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("Set Schema for proxy service", func() {

	var args []string
	var requirementsFactory *testreq.FakeReqFactory

	var ui *testterm.FakeUI
	var serviceRepo *testapi.FakeServiceRepository
	var config core_config.Repository
	var deps command_registry.Dependency

	updateCommandDependency := func(pluginCall bool) {
		deps.Ui = ui
		deps.Config = config
		deps.RepoLocator = deps.RepoLocator.SetServiceRepository(serviceRepo)
		command_registry.Commands.SetCommand(command_registry.Commands.FindCommand("set-schema").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		config = testconfig.NewRepositoryWithDefaults()
	})

	Context("required params", func() {

		It("fails with usage when no args are passed", func() {
			args = []string{}
			requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: true}

			result := testcmd.RunCliCommand("set-schema", args, requirementsFactory, updateCommandDependency, false)

			Expect(result).To(BeFalse())
		})

		It("fails with usage when only one arg is passed", func() {
			args = []string{"my-service"}
			requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: true}

			result := testcmd.RunCliCommand("set-schema", args, requirementsFactory, updateCommandDependency, false)

			Expect(result).To(BeFalse())
		})

		It("fails when not logged in", func() {
			args = []string{"my-service", "schema-file"}
			requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: false, TargetedSpaceSuccess: false}

			result := testcmd.RunCliCommand("set-schema", args, requirementsFactory, updateCommandDependency, false)

			Expect(result).To(BeFalse())
		})

		It("fails when logged in but space is not targetted", func() {
			args = []string{"my-service", "schema-file"}
			requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: false}

			result := testcmd.RunCliCommand("set-schema", args, requirementsFactory, updateCommandDependency, false)

			Expect(result).To(BeFalse())
		})

		It("succeeds when logged in with targetted space and two arguments are passed", func() {
			args = []string{"my-service", "schema-file"}
			requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: true}

			result := testcmd.RunCliCommand("set-schema", args, requirementsFactory, updateCommandDependency, false)

			Expect(result).To(BeTrue())
			Expect(requirementsFactory.ServiceInstanceName).To(Equal("my-service"))
		})

	})

	Context("schema file", func() {

		BeforeEach(func() {
			serviceInstance := models.ServiceInstance{}
			serviceInstance.Name = "my-proxy-service"

			requirementsFactory = &testreq.FakeReqFactory{
				LoginSuccess:         true,
				TargetedSpaceSuccess: true,
				ServiceInstance:      serviceInstance,
			}

			serviceRepo = &testapi.FakeServiceRepository{}
			args = []string{"my-proxy-service", "schemafile.txt"}
		})

		AfterEach(func() {
			os.Remove("schemafile.txt")
		})

		It("fails when schema file does not exists", func() {
			testcmd.RunCliCommand("set-schema", args, requirementsFactory, updateCommandDependency, false)

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Failed to read file schemafile.txt."},
				[]string{"Error: open schemafile.txt: no such file or directory"},
			))
		})

		It("works when empty schema file is given", func() {
			ioutil.WriteFile("schemafile.txt", []byte(nil), 0644)
			testcmd.RunCliCommand("set-schema", args, requirementsFactory, updateCommandDependency, false)

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Applying schema to my-proxy-service in org my-org / space my-space as my-user..."},
				[]string{"OK"},
			))
		})

		It("works when schema file with no duplicates is given", func() {
			ioutil.WriteFile("schemafile.txt", []byte(".sky.com"), 0644)
			testcmd.RunCliCommand("set-schema", args, requirementsFactory, updateCommandDependency, false)

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Applying schema to my-proxy-service in org my-org / space my-space as my-user..."},
				[]string{"OK"},
			))
		})

		It("fails when schema file with duplicates is given", func() {
			ioutil.WriteFile("schemafile.txt", []byte(".sky.com\nupload.sky.com"), 0644)
			testcmd.RunCliCommand("set-schema", args, requirementsFactory, updateCommandDependency, false)

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{`Failed: top level domain 'sky.com' duplicates found: [".sky.com" "upload.sky.com"]`},
			))
		})

	})

})
