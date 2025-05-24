package infra

import (
	"time"
	"wfs/backend/apperr"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type Discord struct {
	url        string
	maxRetry   int
	timeoutSec int
}

func NewDiscord(
	url string,
	maxRetry int,
	timeoutSec int,
) *Discord {
	return &Discord{
		url:        url,
		maxRetry:   maxRetry,
		timeoutSec: timeoutSec,
	}
}

func (d *Discord) Comment(message string) error {
	c := req.C().
		SetBaseURL(d.url).
		SetCommonRetryCount(d.maxRetry).
		SetTimeout(time.Duration(d.timeoutSec) * time.Second)

	_, err := c.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"content": message}).
		Post("")

	if err != nil {
		return failure.Translate(err, apperr.DiscordAPISendLogError)
	}

	return nil
}
