package api

import (
	"fmt"
	"github.com/cloudfoundry/cli/cf/configuration"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/cf/net"
)

type AppSshRepository interface {
	GetSshDetails(appGuid string, instance int) (sshDetails models.SshConnectionDetails, apiErr error)
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

func (repo CloudControllerAppSshRepository) GetSshDetails(appGuid string, instance int) (sshDetails models.SshConnectionDetails, apiErr error) {

	url := fmt.Sprintf("%s/v2/apps/%s/instances/%d/ssh", repo.config.ApiEndpoint(), appGuid, instance)
	request, apiErr := repo.gateway.NewRequest("GET", url, repo.config.AccessToken(), nil)

	if apiErr != nil {
		return
	}

	serverResponse := new(struct {
		Ip     string `json:"ip"`
		Port   int    `json:"port"`
		User   string `json:"user"`
		SshKey string `json:"sshkey"`
	})

	_, apiErr = repo.gateway.PerformRequestForJSONResponse(request, &serverResponse)
	if apiErr != nil {
		return
	}

	sshDetails.Ip = serverResponse.Ip
	sshDetails.Port = serverResponse.Port
	sshDetails.User = serverResponse.User
	sshDetails.SshKey = serverResponse.SshKey

	return
}
