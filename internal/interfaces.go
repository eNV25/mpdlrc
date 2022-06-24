package internal

import (
	"context"
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
	PostEvents(ctx context.Context, ch chan<- tcell.Event)
}

type Song interface {
	ID() string
	Title() string
	Artist() string
	Album() string
	Date() string
	File() string
}

type Status interface {
	Duration() time.Duration
	Elapsed() time.Duration
	State() string
	Repeat() bool
	Random() bool
	Single() bool
	Consume() bool
}

type Widget interface {
	Draw()
	View() views.View
	SetView(view views.View)
	Size() (x int, y int)
	Resize()
}
