package event

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type Draw struct {
	event
}

func NewDrawEvent() tcell.Event {
	return &Draw{event{when: time.Now()}}
}
