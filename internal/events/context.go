package events

import (
	"context"

	"github.com/gdamore/tcell/v2"
)

type key struct{}

func ContextWith(ctx context.Context, events chan<- tcell.Event) context.Context {
	return context.WithValue(ctx, key{}, events)
}

func ContextHas(ctx context.Context) bool {
	return ctx.Value(key{}) != nil
}

func FromContext(ctx context.Context) chan<- tcell.Event {
	return ctx.Value(key{}).(chan<- tcell.Event)
}
