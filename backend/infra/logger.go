package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"wfs/backend/adapter"
	"wfs/backend/config"

	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

type Logger struct {
	zlog   zerolog.Logger
	ownIGN string
}

func NewLogger(i do.Injector) (*Logger, error) {
	config := do.MustInvoke[config.Config](i)
	alertDiscord := do.MustInvokeNamed[adapter.DiscordClient](i, "alert-discord-client")
	infoDiscord := do.MustInvokeNamed[adapter.DiscordClient](i, "info-discord-client")

	zerolog.TimeFieldFormat = time.DateTime
	zerolog.SetGlobalLevel(config.Logger.Level)

	consoleWriter := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}
	remoteWriter := remoteWriter{
		alertDiscord: alertDiscord,
		infoDiscord:  infoDiscord,
	}
	logFile, _ := os.OpenFile(
		filepath.Join(config.LocalFile.RootDir, config.Basic.Name+".log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		os.ModePerm,
	)
	multiLevelWriter := zerolog.MultiLevelWriter(consoleWriter, &remoteWriter, logFile)

	zlog := zerolog.New(multiLevelWriter).
		With().
		Timestamp().
		Str("semver", config.Basic.Version).
		Logger()

	return &Logger{zlog: zlog}, nil
}

func (l *Logger) SetOwnIGN(ownIGN string) {
	l.ownIGN = ownIGN
}

func (l *Logger) Debug(message string, contexts map[string]string) {
	e := l.zlog.Debug().
		Str("ign", l.ownIGN).
		Str("message", message)
	l.addContext(e, contexts)
	e.Send()
}

func (l *Logger) Info(message string, contexts map[string]string) {
	e := l.zlog.Info().
		Str("ign", l.ownIGN).
		Str("message", message)
	l.addContext(e, contexts)
	e.Send()
}

func (l *Logger) Error(err error, contexts map[string]string) {
	e := l.zlog.Error().
		Str("ign", l.ownIGN).
		Str("error", fmt.Sprintf("%+v", err))
	l.addContext(e, contexts)
	e.Send()
}

func (l *Logger) addContext(e *zerolog.Event, contexts map[string]string) {
	if len(contexts) == 0 {
		return
	}

	for key, value := range contexts {
		e = e.Str(key, value)
	}
}

type remoteWriter struct {
	zerolog.FilteredLevelWriter
	alertDiscord adapter.DiscordClient
	infoDiscord  adapter.DiscordClient
}

func (w *remoteWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if level < zerolog.InfoLevel {
		return 0, nil
	}

	var client adapter.DiscordClient
	if level > zerolog.InfoLevel {
		client = w.alertDiscord
	} else {
		client = w.infoDiscord
	}

	formatted := fmt.Sprintf("```%s```", w.pretty(string(p)))

	err := client.Comment(context.Background(), formatted)
	if err != nil {
		log.Printf("Failed to send report: %s\n", err.Error())
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
