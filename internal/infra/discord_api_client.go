package infra

import (
	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type DiscordApiClient interface {
	Comment(message string) error
}

type discordApiClient struct {
	client *req.Client
}

func NewDiscordApiClient(ac apiConfig) DiscordApiClient {
	return &discordApiClient{
		client: req.C().
			SetBaseURL(ac.url).
			SetCommonRetryCount(ac.retryCount).
			SetTimeout(ac.timeout),
	}
}

func (c *discordApiClient) Comment(message string) error {
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
