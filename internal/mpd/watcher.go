package mpd

import (
	"github.com/fhs/gompd/v2/mpd"
	"github.com/gdamore/tcell/v2"

	"github.com/env25/mpdlrc/internal/event"
)

type MPDWatcher struct {
	watcher             *mpd.Watcher
	net, addr, password string
}

func NewMPDWatcher(net, addr, password string) *MPDWatcher {
	return &MPDWatcher{net: net, addr: addr, password: password}
}

func (w *MPDWatcher) Start() (err error) {
	w.watcher, err = mpd.NewWatcher(w.net, w.addr, w.password, "player")
	return
}

func (w *MPDWatcher) Stop() error { return w.watcher.Close() }

func (w *MPDWatcher) PostEvents(postEvent func(tcell.Event) error, quit <-chan struct{}) {
	var newEvent (func() tcell.Event)
	for {
		select {
		case <-quit:
			return
		case mpdev := <-w.watcher.Event:
			switch mpdev {
			case "player":
				newEvent = event.NewPlayerEvent
			}
			if newEvent != nil {
				_ = postEvent(newEvent())
				newEvent = nil
			}
		}
	}
}
