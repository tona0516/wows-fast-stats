package service

import (
	"errors"
	"testing"
	"wfs/backend/application/vo"
	"wfs/backend/domain"

	"github.com/stretchr/testify/assert"
)

func TestUpdater_Updatable_正常系_アップデートあり(t *testing.T) {
	t.Parallel()

	currentVersion := vo.Version{Semver: "1.0.0"}
	mockGithub := &mockGithub{}
	updater := NewUpdater(currentVersion, mockGithub)

	response := domain.GHLatestRelease{TagName: "2.0.0", HTMLURL: "https://hoge.com"}
	mockGithub.On("LatestRelease").Return(response, nil)

	actual, err := updater.Updatable()
	expected := response
	expected.Updatable = true

	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
	mockGithub.AssertCalled(t, "LatestRelease")
}

func TestUpdater_Updatable_正常系_アップデートなし(t *testing.T) {
	t.Parallel()

	currentVersion := vo.Version{Semver: "1.0.0"}
	mockGithub := &mockGithub{}
	updater := NewUpdater(currentVersion, mockGithub)

	response := domain.GHLatestRelease{TagName: "1.0.0", HTMLURL: "https://hoge.com"}
	mockGithub.On("LatestRelease").Return(response, nil)

	actual, err := updater.Updatable()
	expected := response
	expected.Updatable = false

	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
	mockGithub.AssertCalled(t, "LatestRelease")
}

func TestUpdater_Updatable_異常系(t *testing.T) {
	t.Parallel()

	currentVersion := vo.Version{Semver: "1.0.0"}
	mockGithub := &mockGithub{}
	updater := NewUpdater(currentVersion, mockGithub)

	//nolint:goerr113
	expected := errors.New("some error")
	mockGithub.On("LatestRelease").Return(domain.GHLatestRelease{}, expected)

	_, err := updater.Updatable()

	assert.EqualError(t, err, expected.Error())
	mockGithub.AssertCalled(t, "LatestRelease")
}
