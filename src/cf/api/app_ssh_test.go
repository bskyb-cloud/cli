package api_test

import (
	. "cf/api"
	"cf/net"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	testapi "testhelpers/api"
	testconfig "testhelpers/configuration"
	testnet "testhelpers/net"
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

		apiResponse, sshDetails := repo.GetSshDetails("my-app-guid", 0)

		Expect(handler.AllRequestsCalled()).To(BeTrue())
		Expect(apiResponse.IsSuccessful()).To(BeTrue())

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
	ts, handler = testnet.NewTLSServer(requests)
	configRepo := testconfig.NewRepositoryWithDefaults()
	configRepo.SetApiEndpoint(ts.URL)
	gateway := net.NewCloudControllerGateway()
	repo = NewCloudControllerAppSshRepository(configRepo, gateway)
	return
}
