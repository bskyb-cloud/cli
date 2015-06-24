package service_test

import (
	testapi "github.com/cloudfoundry/cli/cf/api/fakes"
	. "github.com/cloudfoundry/cli/cf/commands/service"
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

var (
	args        []string
	reqFactory  *testreq.FakeReqFactory
	serviceRepo *testapi.FakeServiceRepo
)

var _ = Describe("Set Schema for proxy service", func() {

	Context("required params", func() {

		It("fails with usage when no args are passed", func() {
			args = []string{}
			serviceRepo = &testapi.FakeServiceRepo{}
			reqFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: true}

			result, ui := callSetSchema(args, reqFactory, serviceRepo)
			Expect(ui.FailedWithUsage).To(BeTrue())
			Expect(result).To(BeFalse())
		})

		It("fails with usage when only one arg is passed", func() {
			args = []string{"my-service"}
			serviceRepo = &testapi.FakeServiceRepo{}
			reqFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: true}

			result, ui := callSetSchema(args, reqFactory, serviceRepo)
			Expect(ui.FailedWithUsage).To(BeTrue())
			Expect(result).To(BeFalse())
		})

		It("fails when not logged in", func() {
			args = []string{"my-service", "schema-file"}
			serviceRepo = &testapi.FakeServiceRepo{}
			reqFactory = &testreq.FakeReqFactory{LoginSuccess: false, TargetedSpaceSuccess: false}

			result, _ := callSetSchema(args, reqFactory, serviceRepo)
			Expect(result).To(BeFalse())
		})

		It("fails when logged in but space is not targetted", func() {
			args = []string{"my-service", "schema-file"}
			serviceRepo = &testapi.FakeServiceRepo{}
			reqFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: false}

			result, _ := callSetSchema(args, reqFactory, serviceRepo)
			Expect(result).To(BeFalse())
		})

		It("succeeds when logged in with targetted space and two arguments are passed", func() {
			args = []string{"my-service", "schema-file"}
			serviceRepo = &testapi.FakeServiceRepo{}
			reqFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: true}

			result, _ := callSetSchema(args, reqFactory, serviceRepo)
			Expect(result).To(BeTrue())
			Expect(reqFactory.ServiceInstanceName).To(Equal("my-service"))
		})

	})

	Context("schema file", func() {

		BeforeEach(func() {
			serviceInstance := models.ServiceInstance{}
			serviceInstance.Name = "my-proxy-service"

			reqFactory = &testreq.FakeReqFactory{
				LoginSuccess:         true,
				TargetedSpaceSuccess: true,
				ServiceInstance:      serviceInstance,
			}

			serviceRepo = &testapi.FakeServiceRepo{}
			args = []string{"my-proxy-service", "schemafile.txt"}
		})

		AfterEach(func() {
			os.Remove("schemafile.txt")
		})

		It("fails when schema file does not exists", func() {
			_, ui := callSetSchema(args, reqFactory, serviceRepo)

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Failed to read file schemafile.txt."},
				[]string{"Error: open schemafile.txt: no such file or directory"},
			))
		})

		It("works when empty schema file is given", func() {
			ioutil.WriteFile("schemafile.txt", []byte(nil), 0644)
			_, ui := callSetSchema(args, reqFactory, serviceRepo)

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Applying schema to my-proxy-service in org my-org / space my-space as my-user..."},
				[]string{"OK"},
			))
		})

		It("works when schema file with no duplicates is given", func() {
			ioutil.WriteFile("schemafile.txt", []byte(".sky.com"), 0644)
			_, ui := callSetSchema(args, reqFactory, serviceRepo)

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Applying schema to my-proxy-service in org my-org / space my-space as my-user..."},
				[]string{"OK"},
			))
		})

		It("fails when schema file with duplicates is given", func() {
			ioutil.WriteFile("schemafile.txt", []byte(".sky.com\nupload.sky.com"), 0644)
			_, ui := callSetSchema(args, reqFactory, serviceRepo)

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{`Failed: top level domain 'sky.com' duplicates found: [".sky.com" "upload.sky.com"]`},
			))
		})

	})

})

func callSetSchema(args []string, reqFactory *testreq.FakeReqFactory, serviceRepo *testapi.FakeServiceRepo) (bool, *testterm.FakeUI) {
	ui := &testterm.FakeUI{}

	configRepo := testconfig.NewRepositoryWithDefaults()
	cmd := NewSetSchema(ui, configRepo, serviceRepo)
	return testcmd.RunCommand(cmd, args, reqFactory), ui
}
