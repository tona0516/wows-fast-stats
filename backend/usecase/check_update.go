package usecase

import (
	"context"
	"wfs/backend/adapter"
	"wfs/backend/config"
	"wfs/backend/core"

	"github.com/Masterminds/semver/v3"
	"github.com/samber/do/v2"
)

type CheckUpdate struct {
	currentVersion string
	githubClient   adapter.GithubClient
}

func NewCheckUpdate(i do.Injector) (*CheckUpdate, error) {
	config := do.MustInvoke[config.Config](i)
	return &CheckUpdate{
		currentVersion: config.Basic.Version,
		githubClient:   do.MustInvoke[adapter.GithubClient](i),
	}, nil
}

func (c *CheckUpdate) Invoke(ctx context.Context) *core.NewVersion {
	constraint, err := semver.NewConstraint("> " + c.currentVersion)
	if err != nil {
		return nil
	}

	latestRelease, err := c.githubClient.LatestRelease(ctx)
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

	return &core.NewVersion{
		Version:     latestRelease.TagName,
		DownloadURL: latestRelease.HTMLURL,
	}
}
