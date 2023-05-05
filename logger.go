package main

import (
	"changeme/backend/vo"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

const LOG_DIRECTORY = "log"

type Logger struct {
	zlog zerolog.Logger
}

func NewLogger(version vo.Version) *Logger {
	_ = os.Mkdir(LOG_DIRECTORY, 0755)

	writer := &lumberjack.Logger{
		Filename: filepath.Join(LOG_DIRECTORY, "app.log"),
		MaxAge:   14,
		Compress: true,
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
		zlog: zerolog.New(writer).
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
