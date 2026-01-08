package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"wfs/internal/gateway"

	"github.com/rs/zerolog"
)

type remoteWriter struct {
	zerolog.FilteredLevelWriter
	alertDiscord gateway.DiscordClient
	infoDiscord  gateway.DiscordClient
}

func (w *remoteWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if level < zerolog.InfoLevel {
		return 0, nil
	}

	formatted := fmt.Sprintf("```%s```", w.pretty(string(p)))

	var client gateway.DiscordClient
	if level > zerolog.InfoLevel {
		client = w.alertDiscord
	} else {
		client = w.infoDiscord
	}
	err := client.Comment(formatted)
	if err != nil {
		fmt.Printf("Failed to send to discord: %s\n", err.Error())
	}

	return len(p), nil
}

func (w *remoteWriter) pretty(str string) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return str
	}
	return prettyJSON.String()
}
