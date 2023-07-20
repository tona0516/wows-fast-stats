package service

import (
	"context"
	"fmt"
	"time"
	"wfs/backend/application/vo"

	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	EventLog = "LOG"
)

//go:generate stringer -type=LogLevel
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	env        vo.Env
	version    vo.Version
	eventsEmit vo.EventEmit
	report     Report

	ctx context.Context
}

func NewLogger(
	env vo.Env,
	version vo.Version,
	eventsEmit vo.EventEmit,
	report Report,
) *Logger {
	return &Logger{env: env, version: version, eventsEmit: eventsEmit, report: report}
}

func (l *Logger) SetContext(ctx context.Context) {
	l.ctx = ctx
	runtime.LogSetLogLevel(ctx, logger.DEBUG)
}

func (l *Logger) Debug(msg string) {
	runtime.LogDebug(l.ctx, msg)
	l.emitEvent(msg, nil, DEBUG)
}

func (l *Logger) Info(msg string) {
	runtime.LogInfo(l.ctx, msg)
	l.emitEvent(msg, nil, INFO)
}

func (l *Logger) Warn(msg string, err error) {
	runtime.LogWarningf(l.ctx, "%s: %+v", msg, err)
	l.emitEvent(msg, err, WARN)
}

func (l *Logger) Error(msg string, err error) {
	runtime.LogErrorf(l.ctx, "%s: %+v", msg, err)
	l.emitEvent(msg, err, ERROR)
}

func (l *Logger) emitEvent(msg string, err error, logLevel LogLevel) {
	l.eventsEmit(l.ctx, EventLog, vo.LogParam{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		LogLevel:  logLevel.String(),
		Semver:    l.version.Semver,
		Message:   msg,
		Error:     fmt.Sprintf("%+v", err),
	})
}
