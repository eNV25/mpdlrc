// Package mpdconf has a MPD config file scanner.
package mpdconf

import (
	"io"
	"strings"
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
		return unquote(s.s.TokenText())
	}
	return ""
}

// Next advances to the next token. It returns false is there are none left.
func (s *Scanner) Next() bool {
	for {
		tok := s.s.Scan()
		switch tok {
		case scanner.EOF:
			return false
		case '\n':
			continue
		case '#':
			s.skipComment()
			continue
		default:
			s.tok = tok
			return true
		}
	}
}

func (s *Scanner) skipComment() {
	for {
		tok := s.s.Scan()
		switch tok {
		case scanner.EOF, '\n':
			return
		default:
			continue
		}
	}
}

func unquote(s string) string {
	if len(s) <= 2 {
		return ""
	}
	s = strings.ReplaceAll(s, `\"`, `"`)
	return s[1 : len(s)-1]
}
