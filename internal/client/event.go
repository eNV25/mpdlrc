package client

import (
	"github.com/gdamore/tcell/v2"

	"github.com/env25/mpdlrc/internal/event"
)

type PlayerEvent struct {
	event.Event
}

func newPlayerEvent() tcell.Event {
	ev := new(PlayerEvent)
	ev.Init()
	return ev
}

type OptionsEvent struct {
	event.Event
}

func newOptionsEvent() tcell.Event {
	ev := new(OptionsEvent)
	ev.Init()
	return ev
}
