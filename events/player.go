package events

import "github.com/gdamore/tcell/v2"

type PlayerEvent struct {
	*tcell.EventTime
}

func NewPlayerEvent() *PlayerEvent {
	ev := new(tcell.EventTime)
	ev.SetEventNow()
	return &PlayerEvent{ev}
}
