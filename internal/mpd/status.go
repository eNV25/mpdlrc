package mpd

import (
	"strconv"
	"time"

	"github.com/env25/mpdlrc/internal/state"
)

type Status map[string]string

func (s Status) Duration() time.Duration {
	return secondStringToDuration(s["duration"])
}

func (s Status) Elapsed() time.Duration {
	return secondStringToDuration(s["elapsed"])
}

func secondStringToDuration(str string) time.Duration {
	parsed, _ := strconv.ParseFloat(str, 64)
	return time.Duration(parsed * float64(time.Second))
}

func (s Status) State() state.State {
	switch s["state"] {
	case "play":
		return state.Play
	case "stop":
		return state.Stop
	case "pause":
		return state.Pause
	}
	return 0
}
