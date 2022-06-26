package event

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

var (
	_ tcell.Event = event{}
	_ tcell.Event = &event{}
)

type event time.Time

func (ev event) When() time.Time { return time.Time(ev) }

type (
	Draw     struct{ event }
	Ping     struct{ event }
	Player   struct{ event }
	Function struct {
		event
		Func func()
	}
)

func newEvent() event                   { return event(time.Now()) }
func NewDraw() tcell.Event              { return &Draw{newEvent()} }
func NewPing() tcell.Event              { return &Ping{newEvent()} }
func NewPlayer() tcell.Event            { return &Player{newEvent()} }
func NewFunction(fn func()) tcell.Event { return &Function{newEvent(), fn} }
