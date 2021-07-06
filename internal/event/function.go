package event

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type Function struct {
	event

	Run func()
}

func NewFunctionEvent(fn func()) tcell.Event {
	return &Function{event{when: time.Now()}, fn}
}
