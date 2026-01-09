package usecase

import (
	"wfs/internal/adapter"
	"wfs/internal/config"
	"wfs/internal/data"

	"github.com/Masterminds/semver/v3"
	"github.com/samber/do/v2"
)

type UpdateCheck struct {
	currentVersion string
	githubClient   adapter.GithubClient
}

func NewUpdateCheck(i do.Injector) (*UpdateCheck, error) {
	config := do.MustInvoke[config.Config](i)
	return &UpdateCheck{
		currentVersion: config.Basic.Version,
		githubClient:   do.MustInvoke[adapter.GithubClient](i),
	}, nil
}

func (c *UpdateCheck) Invoke() *data.NewVersion {
	constraint, err := semver.NewConstraint("> " + c.currentVersion)
	if err != nil {
		return nil
	}

	latestRelease, err := c.githubClient.LatestRelease()
	if err != nil {
		return nil
	}
	semverVersion, err := semver.NewVersion(latestRelease.TagName)
	if err != nil {
		return nil
	}

	updatable, _ := constraint.Validate(semverVersion)
	if !updatable {
		return nil
	}

	return &data.NewVersion{
		Version:     latestRelease.TagName,
		DownloadURL: latestRelease.HTMLURL,
	}
}
