package events

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type PingEvent struct {
	event
}

func NewPingEvent() tcell.Event {
	return &PingEvent{event{when: time.Now()}}
}
