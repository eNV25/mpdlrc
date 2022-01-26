package lrc

import (
	"reflect"
	"testing"
	"time"
)

var tests = []struct {
	lrc   string
	times []time.Duration
	lines []string
}{
	{
		`
[00:12.00]Line 1 lyrics
[00:17.20]Line 2 lyrics
[00:21.10]Line 3 lyrics
`,
		[]time.Duration{
			parseDuration("00m12.00s"),
			parseDuration("00m17.20s"),
			parseDuration("00m21.10s"),
		},
		[]string{
			"Line 1 lyrics",
			"Line 2 lyrics",
			"Line 3 lyrics",
		},
	},
	{
		`
[00:12.00][00:13.00][00:14.00]Line 1 lyrics
[00:17.20]Line 2 lyrics
[00:21.10][00:22.00]Line 3 lyrics
`,
		[]time.Duration{
			parseDuration("00m12.00s"),
			parseDuration("00m13.00s"),
			parseDuration("00m14.00s"),
			parseDuration("00m17.20s"),
			parseDuration("00m21.10s"),
			parseDuration("00m22.00s"),
		},
		[]string{
			"Line 1 lyrics",
			"Line 1 lyrics",
			"Line 1 lyrics",
			"Line 2 lyrics",
			"Line 3 lyrics",
			"Line 3 lyrics",
		},
	},
}

func parseDuration(text string) (du time.Duration) {
	du, _ = time.ParseDuration(text)
	return
}

func TestParseString(t *testing.T) {
	for i := range tests {
		times, lines, err := ParseString(tests[i].lrc)
		if err != nil || !reflect.DeepEqual(times, tests[i].times) || !reflect.DeepEqual(lines, tests[i].lines) {
			t.Errorf("ParseString(%q) != %v, %v, got = %v, %v", tests[i].lrc, tests[i].times, tests[i].lines, times, lines)
		}
	}
}

func BenchmarkParseString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := range tests {
			ParseString(tests[i].lrc)
		}
	}
}
