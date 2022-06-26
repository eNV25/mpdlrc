package lyrics

import "context"

type key struct{}

func ContextWith(ctx context.Context, lyrics *Lyrics) context.Context {
	return context.WithValue(ctx, key{}, lyrics)
}

func ContextHas(ctx context.Context) bool {
	return ctx.Value(key{}) != nil
}

func FromContext(ctx context.Context) *Lyrics {
	return ctx.Value(key{}).(*Lyrics)
}
