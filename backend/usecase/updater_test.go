package usecase

import (
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/domain/mock_repository"
	"wfs/backend/domain/model"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdater_IsUpdatable(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	t.Run("正常系_アップデートあり", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockGithub := mock_repository.NewMockGithubInterface(ctrl)
		response := model.GHLatestRelease{TagName: "2.0.0", HTMLURL: "https://hoge.com"}
		mockGithub.EXPECT().LatestRelease().Return(response, nil)

		env := model.Env{Semver: "1.0.0"}
		updater := NewUpdater(env, mockGithub, nil)

		// テスト
		actual, err := updater.IsUpdatable()
		expected := response
		expected.Updatable = true

		// アサーション
		assert.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("正常系_アップデートなし", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockGithub := mock_repository.NewMockGithubInterface(ctrl)
		response := model.GHLatestRelease{TagName: "1.0.0", HTMLURL: "https://hoge.com"}
		mockGithub.EXPECT().LatestRelease().Return(response, nil)

		env := model.Env{Semver: "1.0.0"}
		updater := NewUpdater(env, mockGithub, nil)

		// テスト
		actual, err := updater.IsUpdatable()
		expected := response
		expected.Updatable = false

		// アサーション
		assert.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockGithub := mock_repository.NewMockGithubInterface(ctrl)
		expected := failure.New(apperr.HTTPRequestError)
		mockGithub.EXPECT().LatestRelease().Return(model.GHLatestRelease{}, expected)

		env := model.Env{Semver: "1.0.0"}
		updater := NewUpdater(env, mockGithub, nil)

		// テスト
		_, err := updater.IsUpdatable()

		// アサーション
		code, ok := failure.CodeOf(err)
		assert.True(t, ok)
		assert.Equal(t, apperr.HTTPRequestError, code)
	})
}
