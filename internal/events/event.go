package events

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
