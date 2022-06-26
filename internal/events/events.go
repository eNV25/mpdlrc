package events

import (
	"context"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/env25/mpdlrc/internal/event"
)

func Post(ctx context.Context, newEvent func() tcell.Event) {
	select {
	case <-ctx.Done():
	case FromContext(ctx) <- newEvent():
	}
}

func PostEveryTick(ctx context.Context, newEvent func() tcell.Event, d time.Duration) {
	events := FromContext(ctx)
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			select {
			case <-ctx.Done():
				return
			case events <- newEvent():
			}
		}
	}
}

func PostFunc(ctx context.Context, f func()) {
	select {
	case <-ctx.Done():
	case FromContext(ctx) <- event.NewFunction(f):
	}
}
