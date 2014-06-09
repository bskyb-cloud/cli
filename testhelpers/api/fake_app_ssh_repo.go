package api

import (
	"github.com/nimbus-cloud/cli/cf/models"
)

type FakeAppSshRepo struct {
	AppGuid  string
	Instance int
	SshDetails models.SshConnectionDetails
}

func (repo *FakeAppSshRepo) GetSshDetails(appGuid string, instance int) (sshDetails models.SshConnectionDetails, apiErr error) {
	repo.AppGuid = appGuid
	repo.Instance = instance

	sshDetails = repo.SshDetails

	return
}
