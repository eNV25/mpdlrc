package events

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type DrawEvent struct {
	event
}

func NewDrawEvent() tcell.Event {
	return &DrawEvent{event{when: time.Now()}}
}
