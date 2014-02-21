package api

import (
	"cf/net"
	"cf/models"
)

type FakeAppSshRepo struct {
	AppGuid  string
	Instance int
	SshDetails models.SshConnectionDetails
}

func (repo *FakeAppSshRepo) GetSshDetails(appGuid string, instance int) (apiResponse net.ApiResponse, sshDetails models.SshConnectionDetails) {
	repo.AppGuid = appGuid
	repo.Instance = instance

	sshDetails = repo.SshDetails

	return
}
