package infra

import (
	"wfs/backend/domain/model"
	"wfs/backend/infra/webapi"

	"github.com/Masterminds/semver/v3"
	"github.com/morikuni/failure"
)

type VersionFetcher struct {
	github webapi.Github
	semver string
}

func NewVersionFetcher(github webapi.Github, semver string) *VersionFetcher {
	return &VersionFetcher{
		github: github,
		semver: semver,
	}
}

func (f *VersionFetcher) Fetch() (model.LatestRelease, error) {
	var latestRelease model.LatestRelease

	constant, err := semver.NewConstraint("> " + f.semver)
	if err != nil {
		return latestRelease, failure.Wrap(err)
	}

	ghLatestRelease, err := f.github.LatestRelease()
	if err != nil {
		return latestRelease, err
	}
	latest, err := semver.NewVersion(ghLatestRelease.TagName)
	if err != nil {
		return latestRelease, failure.Wrap(err)
	}

	updatable, _ := constant.Validate(latest)
	latestRelease = model.LatestRelease{
		Semver:    ghLatestRelease.TagName,
		URL:       ghLatestRelease.HTMLURL,
		Updatable: updatable,
	}

	return latestRelease, nil
}
