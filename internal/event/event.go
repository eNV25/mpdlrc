package event

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

var (
	_ tcell.Event = Event{}
	_ tcell.Event = &Event{}
)

type Event time.Time

func (ev *Event) Init()          { *ev = Event(time.Now()) }
func (ev Event) When() time.Time { return time.Time(ev) }

type Func struct {
	Event
	Func func()
}

func NewFunction(fn func()) tcell.Event {
	ev := &Func{Func: fn}
	ev.Init()
	return ev
}
