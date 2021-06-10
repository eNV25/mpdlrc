package types

import "time"

type Lyrics interface {
	Lines() []string
	Times() []time.Duration
	Search(time.Duration) int
	N() int
}
