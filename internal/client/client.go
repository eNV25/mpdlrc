package client

import (
	"github.com/env25/mpdlrc/internal/song"
	"github.com/env25/mpdlrc/internal/status"
)

type Client interface {
	NowPlaying() song.Song
	Status() status.Status
	Ping()
	Pause()
	Play()
	Start() error
	Stop() error
}
