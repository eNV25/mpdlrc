package status

import (
	"time"

	"github.com/env25/mpdlrc/internal/app/state"
)

type Status interface {
	Duration() time.Duration
	Elapsed() time.Duration
	State() state.State
}
