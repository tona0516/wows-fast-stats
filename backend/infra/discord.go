package infra

import (
	"net/http"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/infra/webapi"

	"github.com/morikuni/failure"
)

type Discord struct {
	config RequestConfig
}

type DiscordRequestBody struct {
	Content string `json:"content"`
}

func NewDiscord(config RequestConfig) *Discord {
	return &Discord{config: config}
}

func (d *Discord) Comment(message string) error {
	res, err := webapi.PostRequestJSON[DiscordRequestBody, any](
		d.config.URL,
		d.config.Timeout,
		DiscordRequestBody{Content: message},
		d.config.Transport,
	)
	errCtx := failure.Context{
		"url":         res.Request.URL,
		"status_code": strconv.Itoa(res.StatusCode),
		"body":        string(res.BodyByte),
	}
	if err != nil {
		return failure.Wrap(err, errCtx)
	}

	if res.StatusCode != http.StatusOK {
		return failure.New(apperr.DiscordAPISendLogError, errCtx)
	}

	return nil
}
