package usecase

import (
	"context"
	"testing"
	"wfs/backend/adapter"
	"wfs/backend/config"
	"wfs/backend/data"
	"wfs/backend/mock"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateCheck_Invoke(t *testing.T) {
	t.Parallel()

	t.Run("正常系_アップデートあり", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockGithubClient := mock.NewMockGithubClient(ctrl)
		mockGithubClient.EXPECT().LatestRelease(gomock.Any()).Return(data.GHLatestRelease{
			TagName: "2.0.0",
			HTMLURL: "https://hoge.com",
		}, nil)

		injector := do.New()
		do.ProvideValue(injector, config.Config{
			Basic: config.BasicConfig{
				Version: "1.0.0",
			},
		})
		do.Provide(injector, func(i do.Injector) (adapter.GithubClient, error) {
			return mockGithubClient, nil
		})
		do.Provide(injector, NewUpdateCheck)

		uc := do.MustInvoke[*UpdateCheck](injector)
		actual := uc.Invoke(context.Background())

		assert.Equal(t, data.NewVersion{
			Version:     "2.0.0",
			DownloadURL: "https://hoge.com",
		}, *actual)
	})

	t.Run("正常系_アップデートなし", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockGithubClient := mock.NewMockGithubClient(ctrl)
		response := data.GHLatestRelease{TagName: "1.0.0", HTMLURL: "https://hoge.com"}
		mockGithubClient.EXPECT().LatestRelease(gomock.Any()).Return(response, nil)

		injector := do.New()
		do.ProvideValue(injector, config.Config{
			Basic: config.BasicConfig{
				Version: "1.0.0",
			},
		})
		do.Provide(injector, func(i do.Injector) (adapter.GithubClient, error) {
			return mockGithubClient, nil
		})
		do.Provide(injector, NewUpdateCheck)
		uc := do.MustInvoke[*UpdateCheck](injector)
		actual := uc.Invoke(context.Background())

		assert.Nil(t, actual)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockGithubClient := mock.NewMockGithubClient(ctrl)
		expectedErr := failure.New(data.ErrGithubAPI)
		mockGithubClient.EXPECT().LatestRelease(gomock.Any()).Return(data.GHLatestRelease{}, expectedErr)

		injector := do.New()
		do.ProvideValue(injector, config.Config{
			Basic: config.BasicConfig{
				Version: "1.0.0",
			},
		})
		do.Provide(injector, func(i do.Injector) (adapter.GithubClient, error) {
			return mockGithubClient, nil
		})
		do.Provide(injector, NewUpdateCheck)
		uc := do.MustInvoke[*UpdateCheck](injector)
		actual := uc.Invoke(context.Background())

		assert.Nil(t, actual)
	})
}
