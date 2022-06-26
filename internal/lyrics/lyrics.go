package lyrics

import (
	"os"
	"time"

	"github.com/env25/mpdlrc/lrc"
)

type Lyrics struct {
	Times []time.Duration
	Lines []string
}

func New(file string) *Lyrics {
	if r, err := os.Open(file); err == nil {
		if times, lines, err := lrc.ParseReader(r); err == nil {
			return &Lyrics{Times: times, Lines: lines}
		}
	}
	return &Lyrics{}
}
