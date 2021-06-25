package events

import "github.com/gdamore/tcell/v2"

type DrawEvent struct {
	*event
}

func NewDrawEvent() tcell.Event {
	ev := new(event)
	ev.setTimeNow()
	return &DrawEvent{ev}
}
