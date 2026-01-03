package service

import (
	"testing"
	"wfs/internal/apperr"
	"wfs/internal/data"
	"wfs/internal/mock"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateChecker_Invoke(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	t.Run("正常系_アップデートあり", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockGithub := mock.NewMockGithubApiClient(ctrl)
		response := data.GHLatestRelease{TagName: "2.0.0", HTMLURL: "https://hoge.com"}
		mockGithub.EXPECT().LatestRelease().Return(response, nil)

		// テスト
		uc := NewUpdateChecker("1.0.0", mockGithub)
		actual := uc.Invoke()

		// アサーション
		assert.Equal(t, data.NewVersion{
			Semver: "2.0.0",
			URL:    "https://hoge.com",
		}, *actual)
	})

	t.Run("正常系_アップデートなし", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockGithub := mock.NewMockGithubApiClient(ctrl)
		response := data.GHLatestRelease{TagName: "1.0.0", HTMLURL: "https://hoge.com"}
		mockGithub.EXPECT().LatestRelease().Return(response, nil)

		// テスト
		uc := NewUpdateChecker("1.0.0", mockGithub)
		actual := uc.Invoke()

		// アサーション
		assert.Nil(t, actual)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockGithub := mock.NewMockGithubApiClient(ctrl)
		expected := failure.New(apperr.HTTPRequestError)
		mockGithub.EXPECT().LatestRelease().Return(data.GHLatestRelease{}, expected)

		// テスト
		uc := NewUpdateChecker("1.0.0", mockGithub)
		actual := uc.Invoke()

		// アサーション
		assert.Nil(t, actual)
	})
}
