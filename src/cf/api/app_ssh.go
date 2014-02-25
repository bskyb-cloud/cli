package api

import (
	"cf/configuration"
	"cf/models"
	"cf/net"
	"fmt"
)

type AppSshRepository interface {
	GetSshDetails(appGuid string, instance int) (apiResponse net.ApiResponse, sshDetails models.SshConnectionDetails)
}

type CloudControllerAppSshRepository struct {
	config  configuration.Reader
	gateway net.Gateway
}

func NewCloudControllerAppSshRepository(config configuration.Reader, gateway net.Gateway) (repo CloudControllerAppSshRepository) {
	repo.config = config
	repo.gateway = gateway
	return
}

func (repo CloudControllerAppSshRepository) GetSshDetails(appGuid string, instance int) (apiResponse net.ApiResponse, sshDetails models.SshConnectionDetails) {

	url := fmt.Sprintf("%s/v2/apps/%s/instances/%d/ssh", repo.config.ApiEndpoint(), appGuid, instance)
	request, apiResponse := repo.gateway.NewRequest("GET", url, repo.config.AccessToken(), nil)

	if apiResponse.IsNotSuccessful() {
		return
	}

	serverResponse := new(struct {
		Ip     string `json:"ip"`
		Port   int    `json:"port"`
		User   string `json:"user"`
		SshKey string `json:"sshkey"`
	})

	_, apiResponse = repo.gateway.PerformRequestForJSONResponse(request, &serverResponse)
	if apiResponse.IsNotSuccessful() {
		return
	}

	sshDetails.Ip = serverResponse.Ip
	sshDetails.Port = serverResponse.Port
	sshDetails.User = serverResponse.User
	sshDetails.SshKey = serverResponse.SshKey

	return
}
