package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type Logger struct{}

func NewLogger() *Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return &Logger{}
}

func (l *Logger) Info(msg string) {
	log.Info().Msg(msg)
}

func (l *Logger) Warn(msg string, err error) {
	log.Warn().Stack().Err(err).Msg(msg)
}

func (l *Logger) Error(msg string, err error) {
	log.Error().Stack().Err(err).Msg(msg)
}
