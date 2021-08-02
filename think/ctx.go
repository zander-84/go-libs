package think

import "context"

var TraceIDKey = "think_trace_id"

func ContextSetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

func ContextGetTraceID(ctx context.Context) string {
	res, ok := ctx.Value(TraceIDKey).(string)
	if ok {
		return res
	}
	return ""
}
