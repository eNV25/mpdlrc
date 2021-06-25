package events

import "github.com/gdamore/tcell/v2"

type PingEvent struct {
	*event
}

func NewPingEvent() tcell.Event {
	ev := new(event)
	ev.setTimeNow()
	return &PingEvent{ev}
}
