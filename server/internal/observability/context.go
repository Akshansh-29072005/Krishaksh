package observability

import "context"

type key string

const (
	RequestIDKey key = "request_id"
	TraceIDKey   key = "trace_id"
)

func WithRequestID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, RequestIDKey, v)
}
func WithTraceID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, TraceIDKey, v)
}
func RequestIDFrom(ctx context.Context) string {
	v, _ := ctx.Value(RequestIDKey).(string)
	return v
}
func TraceIDFrom(ctx context.Context) string {
	v, _ := ctx.Value(TraceIDKey).(string)
	return v
}
