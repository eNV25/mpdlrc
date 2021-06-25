package events

import "github.com/gdamore/tcell/v2"

type FunctionEvent struct {
	*event
	Run func()
}

func NewFunctionEvent(fn func()) tcell.Event {
	ev := new(event)
	ev.setTimeNow()
	return &FunctionEvent{ev, fn}
}
