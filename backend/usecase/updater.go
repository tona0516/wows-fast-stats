package usecase

import (
	"wfs/backend/domain/model"
	"wfs/backend/domain/repository"

	"github.com/Masterminds/semver/v3"
	"github.com/morikuni/failure"
)

type Updater struct {
	env    model.Env
	github repository.GithubInterface
	logger repository.LoggerInterface
}

func NewUpdater(
	env model.Env,
	github repository.GithubInterface,
	logger repository.LoggerInterface,
) *Updater {
	return &Updater{
		env:    env,
		github: github,
		logger: logger,
	}
}

func (u *Updater) IsUpdatable() (model.GHLatestRelease, error) {
	var latestRelease model.GHLatestRelease

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
