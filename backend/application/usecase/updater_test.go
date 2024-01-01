package usecase

import (
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/application/vo"
	"wfs/backend/domain"
	"wfs/backend/mocks"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdater_IsUpdatable(t *testing.T) {
	t.Parallel()
	t.Run("正常系_アップデートあり", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockGithub := &mocks.GithubInterface{}
		response := domain.GHLatestRelease{TagName: "2.0.0", HTMLURL: "https://hoge.com"}
		mockGithub.On("LatestRelease").Return(response, nil)

		env := vo.Env{Semver: "1.0.0"}
		updater := NewUpdater(env, mockGithub)

		// テスト
		actual, err := updater.IsUpdatable()
		expected := response
		expected.Updatable = true

		// アサーション
		assert.Equal(t, expected, actual)
		require.NoError(t, err)
		mockGithub.AssertExpectations(t)
	})

	t.Run("正常系_アップデートなし", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockGithub := &mocks.GithubInterface{}
		response := domain.GHLatestRelease{TagName: "1.0.0", HTMLURL: "https://hoge.com"}
		mockGithub.On("LatestRelease").Return(response, nil)

		env := vo.Env{Semver: "1.0.0"}
		updater := NewUpdater(env, mockGithub)

		// テスト
		actual, err := updater.IsUpdatable()
		expected := response
		expected.Updatable = false

		// アサーション
		assert.Equal(t, expected, actual)
		require.NoError(t, err)
		mockGithub.AssertExpectations(t)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockGithub := &mocks.GithubInterface{}
		expected := failure.New(apperr.HTTPRequestError)
		mockGithub.On("LatestRelease").Return(domain.GHLatestRelease{}, expected)

		env := vo.Env{Semver: "1.0.0"}
		updater := NewUpdater(env, mockGithub)

		// テスト
		_, err := updater.IsUpdatable()

		// アサーション
		code, ok := failure.CodeOf(err)
		assert.True(t, ok)
		assert.Equal(t, apperr.HTTPRequestError, code)
		mockGithub.AssertExpectations(t)
	})
}
