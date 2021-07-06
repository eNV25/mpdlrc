package events

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type PlayerEvent struct {
	event
}

func NewPlayerEvent() tcell.Event {
	return &PlayerEvent{event{when: time.Now()}}
}
