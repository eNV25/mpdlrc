package types

import (
	"time"
)

type Client interface {
	NowPlaying() Song
	Elapsed() time.Duration
	State() State
	Pause()
	TogglePlay()
	Play()
	Start()
	Stop()
}
