package infra

import (
	"wfs/backend/domain/model"

	"github.com/Masterminds/semver/v3"
	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type VersionFetcher struct {
	githubClient req.Client
	semver       string
}

func NewVersionFetcher(githubClient req.Client, semver string) *VersionFetcher {
	return &VersionFetcher{
		githubClient: githubClient,
		semver:       semver,
	}
}

func (f *VersionFetcher) Fetch() (model.LatestRelease, error) {
	var latestRelease model.LatestRelease

	constant, err := semver.NewConstraint("> " + f.semver)
	if err != nil {
		return latestRelease, failure.Wrap(err)
	}

	//nolint:godox
	// TODO: https://docs.github.com/rest/overview/resources-in-the-rest-api#rate-limiting
	var body GHLatestRelease
	_, err = f.githubClient.R().
		SetSuccessResult(&body).
		Get("/repos/tona0516/wows-fast-stats/releases/latest")
	if err != nil {
		return latestRelease, failure.Wrap(err)
	}

	latest, err := semver.NewVersion(body.TagName)
	if err != nil {
		return latestRelease, failure.Wrap(err)
	}

	updatable, _ := constant.Validate(latest)
	latestRelease = model.LatestRelease{
		Semver:    body.TagName,
		URL:       body.HTMLURL,
		Updatable: updatable,
	}

	return latestRelease, nil
}
