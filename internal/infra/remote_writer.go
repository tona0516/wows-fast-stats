package infra

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
)

type remoteWriter struct {
	zerolog.FilteredLevelWriter
	alertDiscord DiscordApiClient
	infoDiscord  DiscordApiClient
}

func (w *remoteWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	formatted := fmt.Sprintf("```%s```", w.pretty(string(p)))

	var client DiscordApiClient
	if level >= zerolog.WarnLevel {
		client = w.alertDiscord
	} else {
		client = w.infoDiscord
	}
	err := client.Comment(formatted)
	if err != nil {
		fmt.Printf("Failed to send to discord: %s", err.Error())
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
