package main

import (
	"changeme/backend/vo"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const LOG_DIRECTORY = "log"

type Logger struct {
	zlog zerolog.Logger
}

func NewLogger(version vo.Version) *Logger {
	_ = os.Mkdir(LOG_DIRECTORY, 0755)

	now := time.Now()
	f, err := os.OpenFile(
		filepath.Join(LOG_DIRECTORY, now.Format("20060102")+".log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		panic(err)
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	return &Logger{
		zlog: zerolog.New(f).
			With().
			Timestamp().
			Stack().
			Str("version", version.Semver).
			Str("revision", version.Revision).
			Logger(),
	}
}

func (l *Logger) Info(msg string) {
	l.zlog.Info().Caller(1).Msg(msg)
}

func (l *Logger) Warn(msg string, err error) {
	l.zlog.Warn().Caller(1).Err(err).Msg(msg)
}

func (l *Logger) Error(msg string, err error) {
	l.zlog.Error().Caller(1).Err(err).Msg(msg)
}
