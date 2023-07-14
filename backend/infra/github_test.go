package infra

import (
	"testing"
	"wfs/backend/vo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGithub_LatestRelease(t *testing.T) {
	t.Parallel()

	expected := vo.GHLatestRelease{
		HTMLURL: "https://github.com/tona0516/wows-fast-stats/releases/tag/1.0.0",
		TagName: "1.0.0",
	}

	github := NewGithub(vo.RequestConfig{})
	mockAPIClient := &mockAPIClient[vo.GHLatestRelease]{}
	mockAPIClient.On("GetRequest", mock.Anything).Return(APIResponse[vo.GHLatestRelease]{Body: expected}, nil)
	github.latestReleaseClient = mockAPIClient

	actual, err := github.LatestRelease()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
