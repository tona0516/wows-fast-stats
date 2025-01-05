package infra

import (
	"wfs/backend/apperr"
	"wfs/backend/infra/webapi"

	"github.com/morikuni/failure"
)

type Discord struct {
	config RequestConfig
}

func NewDiscord(config RequestConfig) *Discord {
	return &Discord{config: config}
}

func (d *Discord) Comment(message string) error {
	_, _, err := webapi.NewClient(d.config.URL,
		webapi.WithTimeout(d.config.Timeout),
		webapi.WithHeaders(map[string]string{"Content-Type": "application/json"}),
		webapi.WithBody(map[string]string{"content": message}),
	).POST()
	if err != nil {
		return failure.Translate(err, apperr.DiscordAPISendLogError)
	}

	return nil
}
