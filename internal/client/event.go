package client

import (
	"github.com/gdamore/tcell/v2"

	"github.com/env25/mpdlrc/internal/event"
)

type _event struct {
	event.Event
	Data
}

// PlayerEvent occurs when the state of the music player changes.
type PlayerEvent struct{ _event }

func newPlayerEvent(data Data) tcell.Event {
	ev := new(PlayerEvent)
	ev.Init()
	ev.Data = data
	return ev
}

// OptionsEvent occurs when an option in the music player changes.
type OptionsEvent struct{ _event }

func newOptionsEvent(data Data) tcell.Event {
	ev := new(OptionsEvent)
	ev.Init()
	ev.Data = data
	return ev
}
