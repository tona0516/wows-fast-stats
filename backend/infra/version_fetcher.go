package infra

import (
	"encoding/json"
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

	// TODO: https://docs.github.com/rest/overview/resources-in-the-rest-api#rate-limiting
	resp, err := f.githubClient.R().Get("/repos/tona0516/wows-fast-stats/releases/latest")
	if err != nil {
		return latestRelease, failure.Wrap(err)
	}

	var ghLatestRelease GHLatestRelease
	if err := json.Unmarshal(resp.Bytes(), &ghLatestRelease); err != nil {
		return latestRelease, failure.Wrap(err)
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
