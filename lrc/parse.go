package lrc

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"time"
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
	var tt, tmpt time.Duration
	var ll, tmpl []byte
	var err error
	var ok bool

	lines := make([]string, 0)
	times := make([]time.Duration, 0)

	// loop line by line until EOF
	for {
		var is int
		ll, err = p.r.ReadSlice('\n')

		ll = bytes.TrimSpace(ll)

		// parse same line until no match
		for {
			tmpt, tmpl, ok = parseLine(ll)

			if !ok {
				break
			}

			ll = tmpl
			tt = tmpt

			i++
			is++
			times = append(times, tt)
		}

		for x := 0; x < is; x++ {
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

func parseLine(line []byte) (tt time.Duration, ll []byte, ok bool) {
	// 0123456789
	// [00:00.00]

	if len(line) < 10 {
		return 0, nil, false
	}

	// len("[00:00.00]") => 10
	// len("00m00.00s") => 9

	// [00:00.00] => 00m00.00s
	tmp := make([]byte, 0, 9)
	tmp = append(tmp, line[1:3]...)
	tmp = append(tmp, 'm')
	tmp = append(tmp, line[4:9]...)
	tmp = append(tmp, 's')

	{
		du, err := time.ParseDuration(string(tmp))
		if err != nil {
			return 0, nil, false
		}

		tt = du
	}

	{
		ll = append(make([]byte, 0, len(line[10:])), line[10:]...)
	}

	return tt, ll, true
}
