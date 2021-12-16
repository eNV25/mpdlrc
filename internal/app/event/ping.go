package event

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type Ping struct {
	event
}

func NewPingEvent() tcell.Event {
	return &Ping{event{when: time.Now()}}
}
