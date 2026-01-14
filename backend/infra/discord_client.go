package infra

import (
	"context"
	"time"
	"wfs/backend/data"

	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type DiscordClient struct {
	client *req.Client
}

func NewDiscordClient(
	url string,
	retryCount int,
	timeout time.Duration,
) (*DiscordClient, error) {
	return &DiscordClient{
		client: req.C().
			SetBaseURL(url).
			SetCommonRetryCount(retryCount).
			SetTimeout(timeout),
	}, nil
}

func (c *DiscordClient) Comment(ctx context.Context, message string) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"content": message}).
		Post("")
	if err != nil {
		return failure.Translate(err, data.ErrDiscordAPI)
	}

	if resp.IsErrorState() {
		return failure.New(data.ErrDiscordAPI, failure.Context{
			"status_code": resp.Status,
			"body":        resp.String(),
		})
	}

	return nil
}
