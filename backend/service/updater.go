package service

import (
	"fmt"
	"wfs/backend/infra"
	"wfs/backend/vo"

	"github.com/Masterminds/semver/v3"
)

type Updater struct {
	currentVersion vo.Version
	github         infra.GithubInterface
}

func NewUpdater(
	currentVersion vo.Version,
	github infra.GithubInterface,
) *Updater {
	return &Updater{
		currentVersion: currentVersion,
		github:         github,
	}
}

func (u *Updater) Updatable() (vo.GHLatestRelease, error) {
	var latestRelease vo.GHLatestRelease

	c, err := semver.NewConstraint(fmt.Sprintf("> %s", u.currentVersion.Semver))
	if err != nil {
		return latestRelease, err
	}

	latestRelease, err = u.github.LatestRelease()
	if err != nil {
		return latestRelease, err
	}
	latest, err := semver.NewVersion(latestRelease.TagName)
	if err != nil {
		return latestRelease, err
	}

	updatable, _ := c.Validate(latest)
	latestRelease.Updatable = updatable

	return latestRelease, nil
}
