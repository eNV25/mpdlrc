package internal

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type ClientInterface interface {
	NowPlaying() Song
	Status() Status
	Ping()
	Pause()
	Play()
	Start() error
	Stop() error
}

type WatcherInterface interface {
	Start() error
	Stop() error
	PostEvents(ch chan<- tcell.Event, quit <-chan struct{})
}

type SongInterface interface {
	ID() string
	Title() string
	Artist() string
	Album() string
	File() string
	LRCFile() string
}

type StatusInterface interface {
	Duration() time.Duration
	Elapsed() time.Duration
	State() State
}

type Widget interface {
	Draw()
	Resize()
	HandleEvent(ev tcell.Event) bool
	SetView(view views.View)
	Size() (x int, y int)
}
