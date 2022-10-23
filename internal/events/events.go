// Package events implements functions dealing with channels of [tcell.Event].
package events

import (
	"context"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/env25/mpdlrc/internal/event"
	"github.com/env25/mpdlrc/internal/panics"
)

// Post sends a newEvent.
func Post(ctx context.Context, newEvent func() tcell.Event) (done bool) {
	return PostEvent(ctx, newEvent())
}

// PostEvent sends ev.
func PostEvent(ctx context.Context, ev tcell.Event) (done bool) {
	defer panics.Handle(ctx)
	events := FromContext(ctx)
	defer func() { _ = recover() }()
	select {
	case <-ctx.Done():
		return false
	case events <- ev:
	}
	return true
}

// PostFunc sends an event that executes fn.
func PostFunc(ctx context.Context, fn func()) (done bool) {
	return PostEvent(ctx, event.NewFunc(fn))
}

// PostFuncTicker sends an event that executes fn every d duration. Should be called as a new goroutine.
func PostFuncTicker(ctx context.Context, fn func(), d time.Duration) {
	defer panics.Handle(ctx)
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !PostEvent(ctx, event.NewFunc(fn)) {
				return
			}
		}
	}
}
