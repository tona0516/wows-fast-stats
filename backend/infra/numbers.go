package infra

import (
	"strconv"
	"wfs/backend/application/vo"
	"wfs/backend/domain"
	"wfs/backend/logger"

	"github.com/cenkalti/backoff/v4"
)

type Numbers struct {
	config RequestConfig
}

func NewNumbers(config RequestConfig) *Numbers {
	return &Numbers{config: config}
}

func (n *Numbers) ExpectedStats() (domain.NSExpectedStats, error) {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), n.config.Retry)
	operation := func() (APIResponse[domain.NSExpectedStats], error) {
		return getRequest[domain.NSExpectedStats](n.config.URL)
	}

	res, err := backoff.RetryWithData(operation, b)
	if err != nil {
		logger.Error(
			err,
			vo.NewPair("url", n.config.URL),
			vo.NewPair("status_code", strconv.Itoa(res.StatusCode)),
			vo.NewPair("response_body", string(res.BodyByte)),
		)
	}

	return res.Body, err
}
