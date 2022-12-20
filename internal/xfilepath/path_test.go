package xfilepath_test

import (
	"reflect"
	"testing"

	"github.com/env25/mpdlrc/internal/xfilepath"
)

func TestTrimExt(t *testing.T) {
	for _, c := range [...]struct {
		in  string
		out string
	}{
		{"test.mp3", "test"},
		{"test.tar.gz", "test.tar"},
	} {
		out := xfilepath.TrimExt(c.in)
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("TrimExt(%q) = %q, want %q", c.in, out, c.out)
		}
		t.Logf("TrimExt(%q) = %q", c.in, out)
	}
}

func TestReplaceExt(t *testing.T) {
	for _, c := range [...]struct {
		in  string
		ext string
		out string
	}{
		{"test.mp3", ".lrc", "test.lrc"},
		{"test.tar.gz", ".zst", "test.tar.zst"},
	} {
		out := xfilepath.ReplaceExt(c.in, c.ext)
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("ReplaceExt(%q, %q) = %q, want %q", c.in, c.ext, out, c.out)
		}
		t.Logf("ReplaceExt(%q, %q) = %q", c.in, c.ext, out)
	}
}
