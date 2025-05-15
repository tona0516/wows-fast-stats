package service

import (
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/mock/repository"

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
		mockGithub := repository.NewMockGithubInterface(ctrl)
		response := data.GHLatestRelease{TagName: "2.0.0", HTMLURL: "https://hoge.com"}
		mockGithub.EXPECT().LatestRelease().Return(response, nil)

		updater := NewUpdater("1.0.0", mockGithub, nil)

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
		mockGithub := repository.NewMockGithubInterface(ctrl)
		response := data.GHLatestRelease{TagName: "1.0.0", HTMLURL: "https://hoge.com"}
		mockGithub.EXPECT().LatestRelease().Return(response, nil)

		updater := NewUpdater("1.0.0", mockGithub, nil)

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
		mockGithub := repository.NewMockGithubInterface(ctrl)
		expected := failure.New(apperr.HTTPRequestError)
		mockGithub.EXPECT().LatestRelease().Return(data.GHLatestRelease{}, expected)

		updater := NewUpdater("1.0.0", mockGithub, nil)

		// テスト
		_, err := updater.IsUpdatable()

		// アサーション
		code, ok := failure.CodeOf(err)
		assert.True(t, ok)
		assert.Equal(t, apperr.HTTPRequestError, code)
	})
}
