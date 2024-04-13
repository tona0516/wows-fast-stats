package service

import (
	"wfs/backend/data"
	"wfs/backend/repository"

	"github.com/Masterminds/semver/v3"
	"github.com/morikuni/failure"
)

type Updater struct {
	env    data.Env
	github repository.GithubInterface
	logger repository.LoggerInterface
}

func NewUpdater(
	env data.Env,
	github repository.GithubInterface,
	logger repository.LoggerInterface,
) *Updater {
	return &Updater{
		env:    env,
		github: github,
		logger: logger,
	}
}

func (u *Updater) IsUpdatable() (data.GHLatestRelease, error) {
	var latestRelease data.GHLatestRelease

	c, err := semver.NewConstraint("> " + u.env.Semver)
	if err != nil {
		return latestRelease, failure.Wrap(err)
	}

	latestRelease, err = u.github.LatestRelease()
	if err != nil {
		return latestRelease, err
	}
	latest, err := semver.NewVersion(latestRelease.TagName)
	if err != nil {
		return latestRelease, failure.Wrap(err)
	}

	updatable, _ := c.Validate(latest)
	latestRelease.Updatable = updatable

	return latestRelease, nil
}
