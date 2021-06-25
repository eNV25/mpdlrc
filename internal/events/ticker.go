package events

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

func PostTickerEvents(postEvent func(tcell.Event) error, t time.Duration, newEvent func() tcell.Event, quit <-chan struct{}) {
	ticker := time.NewTicker(t)
	defer ticker.Stop()
	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			_ = postEvent(newEvent())
		}
	}
}
