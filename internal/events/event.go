package events

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type event struct {
	time time.Time
}

// check
var _ tcell.Event = event{}
var _ tcell.Event = &event{}

func (ev event) When() time.Time {
	return ev.time
}

func (ev *event) setTimeNow() {
	ev.time = time.Now()
}

func (ev *event) setTime(time time.Time) {
	ev.time = time
}
