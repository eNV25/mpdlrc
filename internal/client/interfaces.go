package client

import (
	"context"
	"time"
)

type Client interface {
	NowPlaying() (Song, error)
	Status() (Status, error)
	MusicDir() (string, error)
	Ping() error
	Pause() error
	Play() error
	Start() error
	Stop() error
}

type Watcher interface {
	Start() error
	Stop() error
	PostEvents(ctx context.Context)
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
