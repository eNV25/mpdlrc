package mpdconf_test

import (
	"embed"
	"testing"

	"github.com/env25/mpdlrc/internal/mpdconf"
)

//go:embed mpd.conf
var fs embed.FS

func TestScanner(t *testing.T) {
	var s mpdconf.Scanner
	f, _ := fs.Open("mpd.conf")
	s.Init(f)
	s.Next()
	const expected = "/home/media/Music"
	if fpath := s.Str("music_directory"); fpath != expected {
		t.Fatalf("music_directory got %q, should be %q", fpath, expected)
	}
}
