package service

import (
	"fmt"
	"wfs/backend/application/repository"
	"wfs/backend/application/vo"
	"wfs/backend/domain"

	"github.com/Masterminds/semver/v3"
)

type Updater struct {
	env    vo.Env
	github repository.GithubInterface
}

func NewUpdater(
	env vo.Env,
	github repository.GithubInterface,
) *Updater {
	return &Updater{
		env:    env,
		github: github,
	}
}

func (u *Updater) Updatable() (domain.GHLatestRelease, error) {
	var latestRelease domain.GHLatestRelease

	c, err := semver.NewConstraint(fmt.Sprintf("> %s", u.env.Semver))
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
