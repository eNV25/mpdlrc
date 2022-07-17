package client

import (
	"github.com/gdamore/tcell/v2"

	"github.com/env25/mpdlrc/internal/event"
)

type _event struct {
	event.Event
	Data
}

type PlayerEvent struct{ _event }

func newPlayerEvent(data Data) tcell.Event {
	ev := new(PlayerEvent)
	ev.Init()
	ev.Data = data
	return ev
}

type OptionsEvent struct{ _event }

func newOptionsEvent(data Data) tcell.Event {
	ev := new(OptionsEvent)
	ev.Init()
	ev.Data = data
	return ev
}
