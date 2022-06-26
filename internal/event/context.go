package event

import (
	"context"

	"github.com/gdamore/tcell/v2"
)

type key struct{}

func ContextWith(ctx context.Context, ev tcell.Event) context.Context {
	return context.WithValue(ctx, key{}, ev)
}

func ContextHas(ctx context.Context) bool {
	return ctx.Value(key{}) != nil
}

func FromContext(ctx context.Context) tcell.Event {
	return ctx.Value(key{}).(tcell.Event)
}
