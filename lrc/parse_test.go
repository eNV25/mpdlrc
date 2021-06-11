package lrc

import (
	"reflect"
	"testing"
	"time"
)

var testCases = []struct {
	lrc    string
	lyrics *Lyrics
}{
	{
		`
[00:12.00]Line 1 lyrics
[00:17.20]Line 2 lyrics
[00:21.10]Line 3 lyrics
		`,
		&Lyrics{
			i: 3,
			times: []time.Duration{
				parseDuration("00m12.00s"),
				parseDuration("00m17.20s"),
				parseDuration("00m21.10s"),
			},
			lines: []string{
				"Line 1 lyrics",
				"Line 2 lyrics",
				"Line 3 lyrics",
			},
		},
	},
	{
		`
[00:12.00][00:13.00][00:14.00]Line 1 lyrics
[00:17.20]Line 2 lyrics
[00:21.10][00:22.00]Line 3 lyrics
		`,
		&Lyrics{
			i: 6,
			times: []time.Duration{
				parseDuration("00m12.00s"),
				parseDuration("00m13.00s"),
				parseDuration("00m14.00s"),
				parseDuration("00m17.20s"),
				parseDuration("00m21.10s"),
				parseDuration("00m22.00s"),
			},
			lines: []string{
				"Line 1 lyrics",
				"Line 1 lyrics",
				"Line 1 lyrics",
				"Line 2 lyrics",
				"Line 3 lyrics",
				"Line 3 lyrics",
			},
		},
	},
}

func parseDuration(text string) (du time.Duration) {
	du, _ = time.ParseDuration(text)
	return
}

func TestParseString(t *testing.T) {
	for _, cs := range testCases {
		l, err := ParseString(cs.lrc)
		if err != nil || !reflect.DeepEqual(l, cs.lyrics) {
			t.Errorf("ParseString(%q) != %v, got = %v", cs.lrc, cs.lyrics, l)
		}
	}
}
