package infra

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"wfs/backend/apperr"

	"github.com/go-resty/resty/v2"
	"github.com/morikuni/failure"
)

type Discord struct {
	baseURL string
}

func NewDiscord(baseURL string) *Discord {
	return &Discord{baseURL: baseURL}
}

func (d *Discord) Comment(message string) error {
	client := resty.New().
		SetTimeout(5 * time.Second)

	resp, err := client.R().
		SetBody(fmt.Sprintf(`{"content": "%s"}`, message)).
		Post(d.baseURL)

	errCtx := failure.Context{
		"url":         d.baseURL,
		"status_code": strconv.Itoa(resp.StatusCode()),
		"body":        string(resp.Body()),
	}
	if err != nil {
		return failure.Wrap(err, errCtx)
	}

	if resp.StatusCode() != http.StatusOK {
		return failure.New(apperr.DiscordAPISendLogError, errCtx)
	}

	return nil
}
