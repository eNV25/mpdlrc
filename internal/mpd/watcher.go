package mpd

import (
	"github.com/env25/mpdlrc/internal/events"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/gdamore/tcell/v2"
)

type MPDWatcher struct {
	watcher   *mpd.Watcher
	net, addr string
	err       error
}

func NewMPDWatcher(net, addr string) *MPDWatcher {
	return &MPDWatcher{net: net, addr: addr}
}

func (w *MPDWatcher) PostEvents(postEvent func(tcell.Event) error, quit <-chan struct{}) {
	for {
		w.watcher, w.err = mpd.NewWatcher(w.net, w.addr, "", "player")
		if w.err == nil {
			break
		}
	}
	defer w.watcher.Close()
	for {
		select {
		case <-quit:
			return
		case v := <-w.watcher.Event:
			switch v {
			case "player":
				_ = postEvent(events.NewPlayerEvent())
			}
		}
	}
}
