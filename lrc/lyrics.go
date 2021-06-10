package lrc

import (
	"sort"
	"time"
)

type Lyrics struct {
	i     int
	times []time.Duration
	lines []string
}

func NewLyrics(times []time.Duration, lines []string) *Lyrics {
	return &Lyrics{
		times: times,
		lines: lines,
		i:     len(lines),
	}
}

func (l *Lyrics) Lines() []string {
	return l.lines
}

func (l *Lyrics) Times() []time.Duration {
	return l.times
}

func (l *Lyrics) N() int {
	return l.i
}

func (l *Lyrics) Search(d time.Duration) int {
	ret := sort.Search(l.i, func(i int) bool { return l.times[i] >= d }) - 1
	if ret < 0 {
		return 0
	} else {
		return ret
	}
}
