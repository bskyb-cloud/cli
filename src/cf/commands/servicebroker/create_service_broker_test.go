package servicebroker_test

import (
	. "cf/commands/servicebroker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
)

var _ = Describe("Testing with ginkgo", func() {
	It("TestCreateServiceBrokerFailsWithUsage", func() {
		reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}
		serviceBrokerRepo := &testapi.FakeServiceBrokerRepo{}

		ui := callCreateServiceBroker([]string{}, reqFactory, serviceBrokerRepo)
		Expect(ui.FailedWithUsage).To(BeTrue())

		ui = callCreateServiceBroker([]string{"1arg"}, reqFactory, serviceBrokerRepo)
		Expect(ui.FailedWithUsage).To(BeTrue())

		ui = callCreateServiceBroker([]string{"1arg", "2arg"}, reqFactory, serviceBrokerRepo)
		Expect(ui.FailedWithUsage).To(BeTrue())

		ui = callCreateServiceBroker([]string{"1arg", "2arg", "3arg"}, reqFactory, serviceBrokerRepo)
		Expect(ui.FailedWithUsage).To(BeTrue())

		ui = callCreateServiceBroker([]string{"1arg", "2arg", "3arg", "4arg"}, reqFactory, serviceBrokerRepo)
		Expect(ui.FailedWithUsage).To(BeFalse())
	})
	It("TestCreateServiceBrokerRequirements", func() {

		reqFactory := &testreq.FakeReqFactory{}
		serviceBrokerRepo := &testapi.FakeServiceBrokerRepo{}
		args := []string{"1arg", "2arg", "3arg", "4arg"}

		reqFactory.LoginSuccess = false
		callCreateServiceBroker(args, reqFactory, serviceBrokerRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())

		reqFactory.LoginSuccess = true
		callCreateServiceBroker(args, reqFactory, serviceBrokerRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeTrue())
	})
	It("TestCreateServiceBroker", func() {

		reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}
		serviceBrokerRepo := &testapi.FakeServiceBrokerRepo{}
		args := []string{"my-broker", "my username", "my password", "http://example.com"}
		ui := callCreateServiceBroker(args, reqFactory, serviceBrokerRepo)

		testassert.SliceContains(ui.Outputs, testassert.Lines{
			{"Creating service broker", "my-broker", "my-user"},
			{"OK"},
		})

		Expect(serviceBrokerRepo.CreateName).To(Equal("my-broker"))
		Expect(serviceBrokerRepo.CreateUrl).To(Equal("http://example.com"))
		Expect(serviceBrokerRepo.CreateUsername).To(Equal("my username"))
		Expect(serviceBrokerRepo.CreatePassword).To(Equal("my password"))
	})
})

func callCreateServiceBroker(args []string, reqFactory *testreq.FakeReqFactory, serviceBrokerRepo *testapi.FakeServiceBrokerRepo) (ui *testterm.FakeUI) {
	ui = &testterm.FakeUI{}
	ctxt := testcmd.NewContext("create-service-broker", args)
	config := testconfig.NewRepositoryWithDefaults()
	cmd := NewCreateServiceBroker(ui, config, serviceBrokerRepo)
	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
