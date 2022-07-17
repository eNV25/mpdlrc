package client

import "context"

type (
	dataKey struct{}
)

func ContextWithData(ctx context.Context, d Data) context.Context {
	return context.WithValue(ctx, dataKey{}, d)
}

func DataFromContext(ctx context.Context) Data {
	return ctx.Value(dataKey{}).(Data)
}
