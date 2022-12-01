package events

import (
	"context"

	"github.com/gdamore/tcell/v2"
)

type key struct{}

// ContextWith returns a [context.Context] with the events channel.
func ContextWith(ctx context.Context, events chan<- tcell.Event) context.Context {
	return context.WithValue(ctx, key{}, events)
}

// FromContext returns the events channel from ctx.
func FromContext(ctx context.Context) chan<- tcell.Event {
	return ctx.Value(key{}).(chan<- tcell.Event)
}
