package service

import (
	"context"
)

type (
	eventEmitFunc func(ctx context.Context, eventName string, optionalData ...any)
)
