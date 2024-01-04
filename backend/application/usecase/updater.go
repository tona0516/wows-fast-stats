package usecase

import (
	"fmt"
	"wfs/backend/adapter"
	"wfs/backend/application/vo"
	"wfs/backend/domain"

	"github.com/Masterminds/semver/v3"
	"github.com/morikuni/failure"
)

type Updater struct {
	env    vo.Env
	github adapter.GithubInterface
	logger adapter.LoggerInterface
}

func NewUpdater(
	env vo.Env,
	github adapter.GithubInterface,
	logger adapter.LoggerInterface,
) *Updater {
	return &Updater{
		env:    env,
		github: github,
		logger: logger,
	}
}

func (u *Updater) IsUpdatable() (domain.GHLatestRelease, error) {
	var latestRelease domain.GHLatestRelease

	c, err := semver.NewConstraint(fmt.Sprintf("> %s", u.env.Semver))
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
