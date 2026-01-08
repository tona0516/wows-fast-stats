package infra

import (
	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

type DiscordClient struct {
	client *req.Client
}

func NewDiscordClient(ac apiConfig) *DiscordClient {
	return &DiscordClient{
		client: req.C().
			SetBaseURL(ac.url).
			SetCommonRetryCount(ac.retryCount).
			SetTimeout(ac.timeout),
	}
}

func (c *DiscordClient) Comment(message string) error {
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"content": message}).
		Post("")
	if err != nil {
		return failure.Wrap(err)
	}

	if resp.IsErrorState() {
		return failure.Wrap(ErrErrorResponse)
	}

	return nil
}
