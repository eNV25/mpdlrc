package mpd

import (
	"github.com/env25/mpdlrc/internal/events"
	"github.com/fhs/gompd/v2/mpd"
	"github.com/gdamore/tcell/v2"
)

type MPDWatcher struct {
	watcher             *mpd.Watcher
	net, addr, password string
}

func NewMPDWatcher(net, addr, password string) *MPDWatcher {
	return &MPDWatcher{net: net, addr: addr, password: password}
}

func (w *MPDWatcher) PostEvents(postEvent func(tcell.Event) error, quit <-chan struct{}) {
	subsystems := []string{"player"}

	w.watcher, _ = mpd.NewWatcher(w.net, w.addr, w.password, subsystems...)
	defer w.watcher.Close()

	var newEvent (func() tcell.Event)

	for {
		select {
		case <-quit:
			return
		case mpdev := <-w.watcher.Event:
			newEvent = nil

			switch mpdev {
			case "player":
				newEvent = events.NewPlayerEvent
			}

			if newEvent != nil {
				_ = postEvent(newEvent())
			}
		}
	}
}
