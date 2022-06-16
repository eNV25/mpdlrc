package internal

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type Client interface {
	NowPlaying() SongType
	Status() StatusType
	Ping()
	Pause()
	Play()
	Start() error
	Stop() error
}

type Watcher interface {
	Start() error
	Stop() error
	PostEvents(ch chan<- tcell.Event, quit <-chan struct{})
}

type Song interface {
	ID() string
	Title() string
	Artist() string
	Album() string
	File() string
	LRCFile() string
}

type Status interface {
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
