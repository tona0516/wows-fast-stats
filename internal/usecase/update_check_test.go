package usecase

import (
	"errors"
	"testing"
	"wfs/internal/data"
	"wfs/internal/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateCheck_Invoke(t *testing.T) {
	t.Parallel()

	t.Run("正常系_アップデートあり", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockGithub := mock.NewMockGithubApiClient(ctrl)
		mockGithub.EXPECT().LatestRelease().Return(data.GHLatestRelease{
			TagName: "2.0.0",
			HTMLURL: "https://hoge.com",
		}, nil)

		uc := NewUpdateCheck("1.0.0", mockGithub)
		actual := uc.Invoke()

		assert.Equal(t, data.NewVersion{
			Version:     "2.0.0",
			DownloadURL: "https://hoge.com",
		}, *actual)
	})

	t.Run("正常系_アップデートなし", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockGithub := mock.NewMockGithubApiClient(ctrl)
		response := data.GHLatestRelease{TagName: "1.0.0", HTMLURL: "https://hoge.com"}
		mockGithub.EXPECT().LatestRelease().Return(response, nil)

		uc := NewUpdateCheck("1.0.0", mockGithub)
		actual := uc.Invoke()

		assert.Nil(t, actual)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockGithub := mock.NewMockGithubApiClient(ctrl)
		mockGithub.EXPECT().LatestRelease().Return(data.GHLatestRelease{}, errors.New("some error"))

		uc := NewUpdateCheck("1.0.0", mockGithub)
		actual := uc.Invoke()

		assert.Nil(t, actual)
	})
}
