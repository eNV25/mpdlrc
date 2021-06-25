package events

import "github.com/gdamore/tcell/v2"

type PlayerEvent struct {
	*event
}

func NewPlayerEvent() tcell.Event {
	ev := new(event)
	ev.setTimeNow()
	return &PlayerEvent{ev}
}
