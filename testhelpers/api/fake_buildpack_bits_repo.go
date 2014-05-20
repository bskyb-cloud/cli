package api

import (
	"github.com/nimbus-cloud/cli/cf/errors"
	"github.com/nimbus-cloud/cli/cf/models"
)

type FakeBuildpackBitsRepository struct {
	UploadBuildpackErr         bool
	UploadBuildpackApiResponse error
	UploadBuildpackPath        string
}

func (repo *FakeBuildpackBitsRepository) UploadBuildpack(buildpack models.Buildpack, dir string) error {
	if repo.UploadBuildpackErr {
		return errors.New("Invalid buildpack")
	}

	repo.UploadBuildpackPath = dir
	return repo.UploadBuildpackApiResponse
}
