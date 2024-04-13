package repository

import (
	"context"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOPACKAGE/$GOFILE -package $GOPACKAGE
type LoggerInterface interface {
	Init(appCtx context.Context)
	SetOwnIGN(ownIGN string)
	Debug(message string, contexts map[string]string)
	Info(message string, contexts map[string]string)
	Warn(err error, contexts map[string]string)
	Error(err error, contexts map[string]string)
	Fatal(err error, contexts map[string]string)
}
