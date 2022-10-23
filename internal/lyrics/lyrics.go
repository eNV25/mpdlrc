// Package lyrics implements the [Lyrics] structure.
package lyrics

import (
	"os"
	"sort"
	"time"

	"github.com/env25/mpdlrc/internal/ufilepath"
	"github.com/env25/mpdlrc/lrc"
)

// Lyrics holds lyrics.
type Lyrics struct {
	Times []time.Duration
	Lines []string
}

func (l *Lyrics) check() {
	if len(l.Times) != len(l.Lines) {
		panic("lyrics: length of Times and Lines are not equal")
	}
}

// Len is the number of lines.
func (l *Lyrics) Len() int {
	return len(l.Times)
}

// Less reports whether the element with index i
// must sort before the element with index j.
func (l *Lyrics) Less(i, j int) bool {
	return l.Times[i] < l.Times[j]
}

// Swap swaps the lines with indexes i and j.
func (l *Lyrics) Swap(i, j int) {
	l.Times[i], l.Times[j], l.Lines[i], l.Lines[j] = l.Times[j], l.Times[i], l.Lines[j], l.Lines[i]
}

// Sort sorts lines in ascending order.
func (l *Lyrics) Sort() {
	l.check()
	sort.Sort(l)
}

// Search return the first index greater with time than x.
// The return value minus one is the last index less than or equal to x.
func (l *Lyrics) Search(x time.Duration) int {
	l.check()
	return sort.Search(l.Len(), func(i int) bool {
		return l.Times[i] > x
	})
}

func newLyrics(file string) (l *Lyrics) {
	if r, err := os.Open(file); err == nil {
		defer r.Close()
		if times, lines, err := lrc.ParseReader(r); err == nil {
			l = &Lyrics{Times: times, Lines: lines}
			l.check()
		}
	}
	return
}

// ForFile returns [Lyrics] the file in disk.
func ForFile(file string) *Lyrics {
	return newLyrics(ufilepath.ReplaceExt(file, ".lrc"))
}
