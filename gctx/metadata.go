package gctx

import (
	"context"
)

type requestIDKeyType struct{}

var RequestIDKey requestIDKeyType = requestIDKeyType{}

func InjectRequestID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, RequestIDKey, rid)
}

func (k requestIDKeyType) String() string {
	return string("X-Request-ID")
}

func RequestID(ctx context.Context) string {
	rid, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return ""
	}
	return rid
}
