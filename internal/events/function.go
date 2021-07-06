package events

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type FunctionEvent struct {
	event

	Run func()
}

func NewFunctionEvent(fn func()) tcell.Event {
	return &FunctionEvent{event{when: time.Now()}, fn}
}
