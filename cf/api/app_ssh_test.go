package api_test

import (
	. "github.com/nimbus-cloud/cli/cf/api"
	"github.com/nimbus-cloud/cli/cf/net"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	testapi "github.com/nimbus-cloud/cli/testhelpers/api"
	testconfig "github.com/nimbus-cloud/cli/testhelpers/configuration"
	testnet "github.com/nimbus-cloud/cli/testhelpers/net"
)

var _ = Describe("AppSshRepository", func() {
	It("TestGetSshCurrentSpace", func() {
		getAppSshInfoRequest := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
			Method:   "GET",
			Path:     "/v2/apps/my-app-guid/instances/0/ssh",
			Response: testnet.TestResponse{Status: http.StatusOK, Body: getSshInfoResponseBody},
		})

		ts, handler, repo := createSshInfoRepo([]testnet.TestRequest{getAppSshInfoRequest})
		defer ts.Close()

		sshDetails, apiErr := repo.GetSshDetails("my-app-guid", 0)

		Expect(handler.AllRequestsCalled()).To(BeTrue())
		Expect(handler).To(testnet.HaveAllRequestsCalled())
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
	configRepo.SetApiEndpoint(ts.URL)
	gateway := net.NewCloudControllerGateway(configRepo)
	repo = NewCloudControllerAppSshRepository(configRepo, gateway)
	return
}
