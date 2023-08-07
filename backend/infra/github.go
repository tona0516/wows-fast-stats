package infra

import (
	"strconv"
	"wfs/backend/application/vo"
	"wfs/backend/domain"
	"wfs/backend/logger"

	"github.com/cenkalti/backoff/v4"
)

type Github struct {
	config RequestConfig
}

func NewGithub(config RequestConfig) *Github {
	return &Github{config: config}
}

func (g *Github) LatestRelease() (domain.GHLatestRelease, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), g.config.Retry)
	operation := func() (APIResponse[domain.GHLatestRelease], error) {
		return getRequest[domain.GHLatestRelease](
			g.config.URL + "/repos/tona0516/wows-fast-stats/releases/latest",
		)
	}

	res, err := backoff.RetryWithData(operation, b)
	if err != nil {
		logger.Error(
			err,
			vo.NewPair("url", g.config.URL),
			vo.NewPair("status_code", strconv.Itoa(res.StatusCode)),
			vo.NewPair("response_body", string(res.BodyByte)),
		)
	}

	return res.Body, err
}
