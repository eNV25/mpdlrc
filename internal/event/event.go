// Package event implements [Event] and derived types.
package event

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

var (
	_ tcell.Event = Event{}
	_ tcell.Event = &Event{}
)

// Event implements [tcell.Event]
type Event time.Time

// Init initializes the event with the current time.
func (ev *Event) Init() { *ev = Event(time.Now()) }

// When reports the time when the event was generated.
func (ev Event) When() time.Time { return time.Time(ev) }

// Func is a [tcell.Event] that holds a function to be executed.
type Func struct {
	Event
	Func func()
}

// NewFunc return a [Func] event with fn.
func NewFunc(fn func()) tcell.Event {
	ev := &Func{Func: fn}
	ev.Init()
	return ev
}
