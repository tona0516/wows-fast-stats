package infra

import (
	"fmt"
	"net/http"
	"strconv"
	"wfs/backend/apperr"

	"github.com/go-resty/resty/v2"
	"github.com/morikuni/failure"
)

type Discord struct {
	config apiConfig
}

func NewDiscord(config apiConfig) *Discord {
	return &Discord{config: config}
}

func (d *Discord) Comment(message string) error {
	client := resty.New().
		SetTimeout(d.config.timeout).
		SetRetryCount(d.config.retryCount)

	resp, err := client.R().
		SetBody(fmt.Sprintf(`{"content": "%s"}`, message)).
		Post(d.config.baseURL)

	errCtx := failure.Context{
		"url":         d.config.baseURL,
		"status_code": strconv.Itoa(resp.StatusCode()),
		"body":        string(resp.Body()),
	}

	if err != nil {
		return failure.New(apperr.DiscordAPISendLogError, failure.Messagef(err.Error()), errCtx)
	}

	if resp.StatusCode() != http.StatusOK {
		return failure.New(apperr.DiscordAPISendLogError, errCtx)
	}

	return nil
}
