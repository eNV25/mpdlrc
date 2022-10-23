package client

import "context"

type (
	dataKey struct{}
)

// ContextWithData returns a subcontext with the data.
func ContextWithData(ctx context.Context, data Data) context.Context {
	return context.WithValue(ctx, dataKey{}, data)
}

// DataFromContext returns the [Data] in ctx.
func DataFromContext(ctx context.Context) Data {
	return ctx.Value(dataKey{}).(Data)
}
