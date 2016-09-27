package api_test

import (
	. "code.cloudfoundry.org/cli/cf/api"
	"code.cloudfoundry.org/cli/cf/api/apifakes"
	"code.cloudfoundry.org/cli/cf/net"
	"code.cloudfoundry.org/cli/cf/terminal/terminalfakes"
	testconfig "code.cloudfoundry.org/cli/testhelpers/configuration"

	"code.cloudfoundry.org/cli/cf/trace/tracefakes"
	. "code.cloudfoundry.org/cli/testhelpers/matchers"
	testnet "code.cloudfoundry.org/cli/testhelpers/net"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("AppSshRepository", func() {
	It("TestGetSshCurrentSpace", func() {
		getAppSshInfoRequest := apifakes.NewCloudControllerTestRequest(testnet.TestRequest{
			Method:   "GET",
			Path:     "/v2/apps/my-app-guid/instances/0/ssh",
			Response: testnet.TestResponse{Status: http.StatusOK, Body: getSshInfoResponseBody},
		})

		ts, handler, repo := createSshInfoRepo([]testnet.TestRequest{getAppSshInfoRequest})
		defer ts.Close()

		sshDetails, apiErr := repo.GetSshDetails("my-app-guid", 0)

		Expect(handler.AllRequestsCalled()).To(BeTrue())
		Expect(handler).To(HaveAllRequestsCalled())
		Expect(apiErr).To(BeNil())

		Expect(sshDetails.Ip).To(Equal("10.0.0.1"))
		Expect(sshDetails.Port).To(Equal(1234))
		Expect(sshDetails.User).To(Equal("vcap"))
		Expect(sshDetails.SshKey).To(Equal("fakekey"))
	})
})

var getSshInfoResponseBody = `
{ 
	"ip": "10.0.0.1",
	"sshkey": "fakekey",
	"user": "vcap",
	"port": 1234
}`

func createSshInfoRepo(requests []testnet.TestRequest) (ts *httptest.Server, handler *testnet.TestHandler, repo AppSshRepository) {
	ts, handler = testnet.NewServer(requests)
	configRepo := testconfig.NewRepositoryWithDefaults()
	configRepo.SetAPIEndpoint(ts.URL)
	gateway := net.NewCloudControllerGateway(configRepo, time.Now, new(terminalfakes.FakeUI), new(tracefakes.FakePrinter), "")
	repo = NewCloudControllerAppSshRepository(configRepo, gateway)
	return
}
