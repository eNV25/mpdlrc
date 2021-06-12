package mpd

import (
	"github.com/env25/mpdlrc/internal/events"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/gdamore/tcell/v2"
)

type MPDWatcher struct {
	watcher *mpd.Watcher
}

func NewMPDWatcher(net, addr string) *MPDWatcher {
	watcher, err := mpd.NewWatcher(net, addr, "", "player")
	if err != nil {
		panic(err)
	}
	return &MPDWatcher{watcher}
}

func (w *MPDWatcher) PostEvents(postEvent func(tcell.Event) error, quit <-chan struct{}) {
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
