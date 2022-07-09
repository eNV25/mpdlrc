package events

import (
	"context"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/env25/mpdlrc/internal/event"
	"github.com/env25/mpdlrc/internal/panics"
)

func Post(ctx context.Context, newEvent func() tcell.Event) {
	PostEvent(ctx, newEvent())
}

func PostFuncTicker(ctx context.Context, f func(), d time.Duration) {
	defer panics.Handle(ctx)
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			PostEvent(ctx, event.NewFunction(f))
		}
	}
}

func PostEvent(ctx context.Context, ev tcell.Event) {
	defer panics.Handle(ctx)
	events := FromContext(ctx)
	defer func() { _ = recover() }()
	select {
	case <-ctx.Done():
		return
	case events <- ev:
	}
}

func PostFunc(ctx context.Context, f func()) {
	PostEvent(ctx, event.NewFunction(f))
}
