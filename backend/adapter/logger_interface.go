package adapter

import (
	"context"
	"wfs/backend/application/vo"
)

type LoggerInterface interface {
	Init(appCtx context.Context)
	Debug(message string, contexts ...vo.Pair)
	Info(message string, contexts ...vo.Pair)
	Warn(err error, contexts ...vo.Pair)
	Error(err error, contexts ...vo.Pair)
}
