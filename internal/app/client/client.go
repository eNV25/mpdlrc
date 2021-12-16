package client

import (
	"github.com/env25/mpdlrc/internal/app/song"
	"github.com/env25/mpdlrc/internal/app/status"
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
