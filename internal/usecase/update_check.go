package usecase

import (
	"wfs/internal/data"
	"wfs/internal/infra"

	"github.com/Masterminds/semver/v3"
)

type UpdateCheck struct {
	currentVersion string
	github         infra.GithubApiClient
}

func NewUpdateCheck(
	currentVersion string,
	github infra.GithubApiClient,
) *UpdateCheck {
	return &UpdateCheck{
		currentVersion: currentVersion,
		github:         github,
	}
}

func (c *UpdateCheck) Invoke() *data.NewVersion {
	constraint, err := semver.NewConstraint("> " + c.currentVersion)
	if err != nil {
		return nil
	}

	latestRelease, err := c.github.LatestRelease()
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
