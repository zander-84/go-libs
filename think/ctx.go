package think

import "context"

var CtxTraceIDKey = "thinkTraceID"

func ContextSetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, CtxTraceIDKey, traceID)
}

func ContextGetTraceID(ctx context.Context) string {
	res, ok := ctx.Value(CtxTraceIDKey).(string)
	if ok {
		return res
	}
	return ""
}
