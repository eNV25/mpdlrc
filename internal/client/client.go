package client

import (
	"time"

	"github.com/env25/mpdlrc/internal/song"
	"github.com/env25/mpdlrc/internal/state"
)

type Client interface {
	NowPlaying() song.Song
	Elapsed() time.Duration
	State() state.State
	Pause()
	TogglePlay()
	Play()
	Start()
	Stop()
}
