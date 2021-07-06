package event

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type Player struct {
	event
}

func NewPlayerEvent() tcell.Event {
	return &Player{event{when: time.Now()}}
}
