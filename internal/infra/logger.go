package infra

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"wfs/internal/config"
	"wfs/internal/gateway"

	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

type Logger struct {
	zlog   zerolog.Logger
	ownIGN string
}

func NewLogger(i do.Injector) (*Logger, error) {
	config := do.MustInvoke[config.Config](i)
	alertDiscord := do.MustInvokeNamed[gateway.DiscordClient](i, "alert-discord-client")
	infoDiscord := do.MustInvokeNamed[gateway.DiscordClient](i, "info-discord-client")

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
		filepath.Join(config.LocalFile.ConfigDir, config.Basic.Name+".log"),
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
