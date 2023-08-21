package service

import (
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/application/vo"
	"wfs/backend/domain"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
)

func TestUpdater_Updatable_正常系_アップデートあり(t *testing.T) {
	t.Parallel()

	env := vo.Env{Semver: "1.0.0"}
	mockGithub := &mockGithub{}
	updater := NewUpdater(env, mockGithub)

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

	env := vo.Env{Semver: "1.0.0"}
	mockGithub := &mockGithub{}
	updater := NewUpdater(env, mockGithub)

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

	env := vo.Env{Semver: "1.0.0"}
	mockGithub := &mockGithub{}
	updater := NewUpdater(env, mockGithub)

	expected := failure.New(apperr.HTTPRequestError)
	mockGithub.On("LatestRelease").Return(domain.GHLatestRelease{}, expected)

	_, err := updater.Updatable()

	code, ok := failure.CodeOf(err)
	assert.True(t, ok)
	assert.Equal(t, code, apperr.HTTPRequestError)

	mockGithub.AssertCalled(t, "LatestRelease")
}
