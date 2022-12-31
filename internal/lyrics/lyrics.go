// Package lyrics implements the [Lyrics] structure.
package lyrics

import (
	"os"
	"sort"
	"time"

	"github.com/env25/mpdlrc/internal/xfilepath"
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

type implLyricsSort Lyrics

func (l *implLyricsSort) Len() int {
	return len(l.Times)
}

func (l *implLyricsSort) Less(i, j int) bool {
	return l.Times[i] < l.Times[j]
}

func (l *implLyricsSort) Swap(i, j int) {
	l.Times[i], l.Times[j], l.Lines[i], l.Lines[j] = l.Times[j], l.Times[i], l.Lines[j], l.Lines[i]
}

var _ sort.Interface = (*implLyricsSort)(nil)

// Sort sorts lines in ascending order.
func (l *Lyrics) Sort() {
	l.check()
	sort.Sort((*implLyricsSort)(l))
}

// Search return the first index greater with time than x.
// The return value minus one is the last index less than or equal to x.
//
// The lyrics must be sorted by time.
func (l *Lyrics) Search(x time.Duration) int {
	l.check()
	return sort.Search(len(l.Times), func(i int) bool {
		return l.Times[i] >= x
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
	return newLyrics(xfilepath.ReplaceExt(file, ".lrc"))
}
