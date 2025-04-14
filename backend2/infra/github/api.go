package github

import (
	"net/http"
	"wfs/backend/apperr"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type API interface {
	FetchLatestRelease() (LatestRelease, error)
}

type api struct {
	client *req.Client
}

func NewGithub(client *req.Client) API {
	return &api{client: client}
}

func (a *api) FetchLatestRelease() (LatestRelease, error) {
	result := LatestRelease{}

	resp, err := a.client.R().
		SetSuccessResult(&result).
		Get("/repos/tona0516/wows-fast-stats/releases/latest")

	if resp.StatusCode != http.StatusOK {
		return result, failure.New(apperr.GithubAPICheckUpdateError)
	}

	return result, failure.Translate(err, apperr.GithubAPICheckUpdateError)
}
