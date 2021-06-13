package lrc

import (
	"sort"
	"time"
)

// Lyrics holds parsed lyrics.
type Lyrics struct {
	i     int
	times []time.Duration
	lines []string
}

// NewLyrics return a new instance of Lyrics from an slice of times and of lines of a song.
// The user needs make sure that both are the same size.
func NewLyrics(times []time.Duration, lines []string) *Lyrics {
	return &Lyrics{
		times: times,
		lines: lines,
		i:     len(lines),
	}
}

// Lines returns lines.
func (l *Lyrics) Lines() []string {
	return l.lines
}

// Times returns times.
func (l *Lyrics) Times() []time.Duration {
	return l.times
}

// N returns the number of lines held.
func (l *Lyrics) N() int {
	return l.i
}

// Search finds the index of lyrics lines where times is larger than d.
//
// This can be used to find the lyrics line to be displayed at a perticular time.
//
//	i = lyrics.Search(1 * time.Minute) - 1
//
func (l *Lyrics) Search(d time.Duration) int {
	// binary search
	return sort.Search(l.i, func(i int) bool { return l.times[i] >= d })
}
