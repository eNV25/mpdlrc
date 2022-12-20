package dirs_test

import (
	"os"
	"os/user"
	"reflect"
	"runtime"
	"testing"

	"github.com/env25/mpdlrc/internal/dirs"
)

func TestExpandTilde(t *testing.T) {
	const s = string(os.PathSeparator)
	ts := []*struct {
		in  string
		out string
	}{
		{"", ""},
		{s, s},
		{"~", dirs.HomeDir("")},
		{"~" + s, dirs.HomeDir("") + s},
		{"~" + s + "file", dirs.HomeDir("") + s + "file"},
		{"file", "file"},
		{"some" + s + "random" + s + "file", "some" + s + "random" + s + "file"},
		{s + "some" + s + "random" + s + "file", s + "some" + s + "random" + s + "file"},
	}
	switch runtime.GOOS {
	case "linux":
		current, _ := user.Current()
		ts = append(ts, []*struct {
			in  string
			out string
		}{
			{"~" + current.Username, dirs.HomeDir(current.Username)},
			{"~" + current.Username + s, dirs.HomeDir(current.Username) + s},
			{"~" + current.Username + s + "file", dirs.HomeDir(current.Username) + s + "file"},
		}...)
	case "windows":
		ts = append(ts, []*struct {
			in  string
			out string
		}{
			{"~/file\\", dirs.HomeDir("") + "\\file\\"},
		}...)
	}
	for _, c := range ts {
		out := dirs.ExpandTilde(c.in)
		t.Logf("%q => %q", c.in, out)
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("ExpandTilde(%q) = %q, want %q", c.in, out, c.out)
		}
	}
}
