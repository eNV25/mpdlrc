package client

import (
	"context"
	"time"

	"github.com/env25/mpdlrc/internal/lyrics"
)

// Client is an interface implemented by all clients.
type Client interface {
	Close() error
	Data() (Data, error)
	MusicDir() (string, error)
	PostEvents(ctx context.Context)
	TogglePause() bool
}

// Data is an aggregate type for [Song] [Status] and [lyrics.Lyrics].
type Data struct {
	Song
	Status
	*lyrics.Lyrics
}

// Song contains data about a playing in the music player.
type Song interface {
	ID() string
	Title() string
	Artist() string
	Album() string
	Date() string
	File() string
}

// Status contains data about status of the music player.
type Status interface {
	Duration() time.Duration
	Elapsed() time.Duration
	State() string
	Repeat() bool
	Random() bool
	Single() bool
	Consume() bool
}
