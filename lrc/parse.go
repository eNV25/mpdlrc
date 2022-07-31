package lrc

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"time"
)

// NOTE: The parser simply ignores all unmatched lines and doen't return any parse errors.

type (
	Duration = time.Duration
	Text     = string
)

// Parse parses a byte slice of LRC lyrics.
func Parse(data []byte) ([]Duration, []Text, error) {
	return ParseReader(bytes.NewReader(data))
}

// ParseString parses a string of LRC lyrics.
func ParseString(text string) ([]Duration, []Text, error) {
	return ParseReader(strings.NewReader(text))
}

func atoi(a byte) byte {
	return a - '0'
}

// Parse parses the reader according to the LRC format.
// https://en.wikipedia.org/wiki/LRC_(file_format)
func ParseReader(reader io.Reader) ([]Duration, []Text, error) {
	times := make([]Duration, 0)
	lines := make([]Text, 0)
	scnnr := bufio.NewScanner(reader)
	// scnnr.Split(bufio.ScanLines)
	for rp := 0; scnnr.Scan(); {
		ll := scnnr.Text()
	match: // [00:00.00][00:00.00]text -> [00:00.00]text -> text
		for {
			switch {
			case len(ll) >= 10 &&
				ll[0] == '[' &&
				'0' <= ll[1] && ll[1] <= '9' && '0' <= ll[2] && ll[2] <= '9' &&
				ll[3] == ':' &&
				'0' <= ll[4] && ll[4] <= '5' && '0' <= ll[5] && ll[5] <= '9' &&
				ll[6] == '.' &&
				'0' <= ll[7] && ll[7] <= '9' && '0' <= ll[8] && ll[8] <= '9' &&
				ll[9] == ']':
				// ll    -> [00:00.00]
				// index -> 0123456789
				times = append(times, 0+
					Duration(10*atoi(ll[1])+atoi(ll[2]))*time.Minute+
					Duration(10*atoi(ll[4])+atoi(ll[5]))*time.Second+
					Duration(10*atoi(ll[7])+atoi(ll[8]))*time.Second/100)
				ll = ll[10:]
			case len(ll) >= 7 &&
				ll[0] == '[' &&
				'0' <= ll[1] && ll[1] <= '9' && '0' <= ll[2] && ll[2] <= '9' &&
				ll[3] == ':' &&
				'0' <= ll[4] && ll[4] <= '5' && '0' <= ll[5] && ll[5] <= '9' &&
				ll[6] == ']':
				// ll    -> [00:00]
				// index -> 0123456
				times = append(times, 0+
					Duration(10*atoi(ll[1])+atoi(ll[2]))*time.Minute+
					Duration(10*atoi(ll[4])+atoi(ll[5]))*time.Second)
				ll = ll[7:]
			default:
				break match
			}
			rp++
		}
		for ; rp > 0; rp-- {
			lines = append(lines, ll)
		}
	}
	if err := scnnr.Err(); err != nil {
		return nil, nil, err
	}
	return times, lines, nil
}
