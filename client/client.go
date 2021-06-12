package client

import (
	"time"

	"local/mpdlrc/song"
	"local/mpdlrc/state"
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
