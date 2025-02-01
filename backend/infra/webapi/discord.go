package webapi

import (
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type Discord interface {
	Comment(message string) error
}

type discord struct {
	config RequestConfig
}

func NewDiscord(config RequestConfig) Discord {
	return &discord{config: config}
}

func (d *discord) Comment(message string) error {
	_, _, err := NewClient(d.config.URL,
		WithTimeout(d.config.Timeout),
		WithHeaders(map[string]string{"Content-Type": "application/json"}),
		WithBody(map[string]string{"content": message}),
	).POST()
	if err != nil {
		return failure.Translate(err, apperr.DiscordAPISendLogError)
	}

	return nil
}
