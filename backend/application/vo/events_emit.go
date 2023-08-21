package vo

import "context"

type EventEmit func(ctx context.Context, eventName string, optionalData ...interface{})
