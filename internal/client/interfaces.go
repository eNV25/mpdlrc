package client

import (
	"context"
	"time"

	"github.com/env25/mpdlrc/internal/lyrics"
)

type Client interface {
	Close() error
	Data() (Data, error)
	MusicDir() (string, error)
	PostEvents(ctx context.Context)
	TogglePause() bool
}

type Data struct {
	Song
	Status
	*lyrics.Lyrics
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
