package types

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type PlayerEvent struct {
	*tcell.EventTime
}

func NewPlayerEvent() *PlayerEvent {
	ev := new(tcell.EventTime)
	ev.SetEventNow()
	return &PlayerEvent{ev}
}

type TickerEvent struct {
	*tcell.EventTime
}

func NewTickerEvent() *TickerEvent {
	ev := new(tcell.EventTime)
	ev.SetEventNow()
	return &TickerEvent{ev}
}

func PostTickerEvents(postEvent func(tcell.Event) error, t time.Duration, quit <-chan struct{}) {
	ticker := time.NewTicker(t)
	defer ticker.Stop()
	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			_ = postEvent(NewTickerEvent())
		}
	}
}
