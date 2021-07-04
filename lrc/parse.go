package lrc

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"time"
	"unicode"
)

// Parse parses a byte slice of LRC lyrics.
func Parse(data []byte) ([]time.Duration, []string, error) {
	return NewParser(bytes.NewReader(data)).Parse()
}

// ParseString parses a string of LRC lyrics.
func ParseString(text string) ([]time.Duration, []string, error) {
	return NewParser(strings.NewReader(text)).Parse()
}

// Perser is a parser type.
type Parser struct {
	r *bufio.Reader
}

// NewParser return a new parser from a reader.
func NewParser(r io.Reader) *Parser {
	return &Parser{r: bufio.NewReader(r)}
}

// Parse parses the reader according to the LRC format.
// https://en.wikipedia.org/wiki/LRC_(file_format)
func (p *Parser) Parse() ([]time.Duration, []string, error) {
	var i int
	var tt, tmpt time.Duration
	var ll, tmpl string
	var err error
	var ok bool

	lines := make([]string, 0)
	times := make([]time.Duration, 0)

	// loop line by line until EOF
	for {
		var rep int

		ll, err = p.r.ReadString('\n')

		ll = strings.TrimRightFunc(ll, unicode.IsSpace)

		// parse same line until no match
		for {
			tmpt, tmpl, ok = parseLine(ll)

			if !ok {
				break
			}

			ll = tmpl
			tt = tmpt

			i++
			rep++
			times = append(times, tt)
		}

		for x := 0; x < rep; x++ {
			lines = append(lines, ll)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, err
		}
	}

	return times, lines, nil
}

func isDigit(b byte) bool { return '0' <= b && b <= '9' }

func parseLine(line string) (time.Duration, string, bool) {
	// 0123456789
	// [00:00.00]
	// 00m00.00s

	// len("[00:00.00]") => 10
	// len("00m00.00s") => 9

	if len(line) < 10 ||
		line[0] != '[' ||
		!isDigit(line[1]) || !isDigit(line[2]) ||
		line[3] != ':' ||
		!isDigit(line[4]) || !isDigit(line[5]) ||
		line[6] != '.' ||
		!isDigit(line[7]) || !isDigit(line[8]) ||
		line[9] != ']' {

		return 0, "", false
	}

	// [00:00.00] => 00m00.00s
	var tmp strings.Builder
	tmp.WriteString(line[1:3])
	tmp.WriteString("m")
	tmp.WriteString(line[4:9])
	tmp.WriteString("s")

	tt, err := time.ParseDuration(tmp.String())
	if err != nil {
		return 0, "", false
	}

	ll := line[10:]

	return tt, ll, true
}
