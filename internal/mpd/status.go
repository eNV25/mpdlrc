package mpd

import (
	"strconv"
	"time"

	"github.com/env25/mpdlrc/internal/state"
)

type Status map[string]string

func (s Status) Duration() time.Duration {
	duration, _ := strconv.ParseFloat(s["duration"], 64)
	return time.Duration(duration * float64(time.Second))
}

func (s Status) Elapsed() time.Duration {
	elapsed, _ := strconv.ParseFloat(s["elapsed"], 64)
	return time.Duration(elapsed * float64(time.Second))
}

func (s Status) State() state.State {
	switch s["status"] {
	case "play":
		return state.PlayState
	case "stop":
		return state.StopState
	case "pause":
		return state.PauseState
	}
	return 0
}
