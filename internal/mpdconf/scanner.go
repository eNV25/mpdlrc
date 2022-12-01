// Package mpdconf has a MPD config file scanner.
package mpdconf

import (
	"io"
	"strconv"
	"text/scanner"
)

// Scanner implements a scanner for the MPD config file.
type Scanner struct {
	s   scanner.Scanner
	tok rune
}

// Init initialises the scanner.
func (s *Scanner) Init(src io.Reader) {
	s.s.Init(src)
	s.s.Whitespace &^= 1 << '\n' // need to handle ourselves
}

// Str returns the string value associated with key.
func (s *Scanner) Str(key string) string {
	if s.tok == scanner.Ident && key == s.s.TokenText() && s.Next() && s.tok == scanner.String {
		ret, _ := strconv.Unquote(s.s.TokenText())
		return ret
	}
	return ""
}

// Next advances to the next token. It returns false is there are none left.
func (s *Scanner) Next() bool {
	tok := s.skipNewlines()
	for {
		switch tok {
		case scanner.EOF:
			return false
		case '#':
			tok = s.skipComment()
			continue
		default:
			s.tok = tok
			return true
		}
	}
}

func (s *Scanner) skipComment() rune {
	for {
		tok := s.s.Scan()
		switch tok {
		case scanner.EOF:
			return tok
		case '\n':
			return s.skipNewlines()
		default:
			continue
		}
	}
}

func (s *Scanner) skipNewlines() rune {
	for {
		tok := s.s.Scan()
		switch tok {
		case '\n':
			continue
		default:
			return tok
		}
	}
}
