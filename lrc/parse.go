package lrc

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
	"time"
)

const (
	PATTERN  = `\[(?P<mm>\d\d)\:(?P<ss>\d\d)\.(?P<xx>\d\d)\](?P<line>.*)`
	PATTERN2 = `\[(?P<time>\d\d\:\d\d\.\d\d)\](?P<line>.*)`
)

var (
	re = regexp.MustCompile(PATTERN2)
)

func Parse(data []byte) (*Lyrics, error) {
	return NewParser(bytes.NewReader(data)).Parse()
}

func ParseString(text string) (*Lyrics, error) {
	return NewParser(strings.NewReader(text)).Parse()
}

type Parser struct {
	r *bufio.Reader
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		r: bufio.NewReader(r),
	}
}

func (p *Parser) Parse() (*Lyrics, error) {
	var i int
	var tt time.Duration
	var ll []byte
	var m [][]byte
	var err error

	lines := make([]string, 0)
	times := make([]time.Duration, 0)

	// loop line by line until EOF
	for {
		ll, err = p.r.ReadSlice('\n')

		// match same line until no match
		for {
			// m = [1: tt, 2: ll]
			m = re.FindSubmatch(ll)

			// if no match; break
			if m == nil {
				break
			}

			// [0-9][0-9]:[0-9][0-9].[0-9][0-9]
			//                =>
			// [0-9][0-9]m[0-9][0-9].[0-9][0-9]s
			m[1] = append(m[1], 's')
			m[1][2] = 'm'

			tt, _ = time.ParseDuration(string(m[1]))

			ll = m[2]

			i++
			times = append(times, tt)
			lines = append(lines, string(ll))
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return &Lyrics{
		i:     i,
		lines: lines,
		times: times,
	}, nil
}
