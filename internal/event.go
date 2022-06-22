package internal

import (
	"context"
	"time"

	"github.com/gdamore/tcell/v2"
)

var (
	_ tcell.Event = event{}
	_ tcell.Event = &event{}
)

type event struct {
	when time.Time
}

func (ev event) When() time.Time {
	return ev.when
}

type (
	EventDraw     struct{ event }
	EventPing     struct{ event }
	EventPlayer   struct{ event }
	EventFunction struct {
		event
		Func func()
	}
)

func NewEventDraw() tcell.Event {
	return &EventDraw{event{when: time.Now()}}
}

func NewEventFunction(fn func()) tcell.Event {
	return &EventFunction{event{when: time.Now()}, fn}
}

func NewEventPing() tcell.Event {
	return &EventPing{event{when: time.Now()}}
}

func NewEventPlayer() tcell.Event {
	return &EventPlayer{event{when: time.Now()}}
}

func sendNewEventEvery(ctx context.Context, ch chan<- tcell.Event, newEvent func() tcell.Event, d time.Duration) {
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
			case ch <- newEvent():
			}
		}
	}
}
