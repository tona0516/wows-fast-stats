package service

import (
	"wfs/internal/data"
	"wfs/internal/infra"

	"github.com/Masterminds/semver/v3"
)

type UpdateChecker struct {
	currentSemver string
	github        infra.GithubApiClient
}

func NewUpdateChecker(
	semver string,
	github infra.GithubApiClient,
) *UpdateChecker {
	return &UpdateChecker{
		currentSemver: semver,
		github:        github,
	}
}

func (c *UpdateChecker) Invoke() *data.NewVersion {
	constraint, err := semver.NewConstraint("> " + c.currentSemver)
	if err != nil {
		return nil
	}

	latestRelease, err := c.github.LatestRelease()
	if err != nil {
		return nil
	}
	latest, err := semver.NewVersion(latestRelease.TagName)
	if err != nil {
		return nil
	}

	updatable, _ := constraint.Validate(latest)
	if !updatable {
		return nil
	}

	return &data.NewVersion{
		Semver: latestRelease.TagName,
		URL:    latestRelease.HTMLURL,
	}
}
