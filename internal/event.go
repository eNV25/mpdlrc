package internal

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

var _ tcell.Event = event{}
var _ tcell.Event = &event{}

type event struct {
	when time.Time
}

func (ev event) When() time.Time {
	return ev.when
}

type EventDraw struct{ event }
type EventPing struct{ event }
type EventPlayer struct{ event }
type EventFunction struct {
	event
	Func func()
}

func NewEventDraw() tcell.Event {
	return &EventDraw{event{when: time.Now()}}
}
func NewEventFunction(fn func()) tcell.Event {
	return &EventFunction{event{when: time.Now()}, fn}
}
func NewEventPing() tcell.Event {
	return &EventPing{event{when: time.Now()}}
}
func NewEventPlayer() tcell.Event {
	return &EventPlayer{event{when: time.Now()}}
}
func sendNewEventEvery(ch chan<- tcell.Event, newEvent func() tcell.Event, d time.Duration, quit <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			ch <- newEvent()
		}
	}
}
