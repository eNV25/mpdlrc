package config

import (
	"os"
	"os/user"
	"reflect"
	"runtime"
	"testing"
)

func TestExpandTilde(t *testing.T) {
	const s = string(os.PathSeparator)
	ts := []*struct {
		in  string
		out string
	}{
		{"", ""},
		{s, s},
		{"~", HomeDir("")},
		{"~" + s, HomeDir("") + s},
		{"~" + s + "file", HomeDir("") + s + "file"},
		{"file", "file"},
		{"some" + s + "random" + s + "file", "some" + s + "random" + s + "file"},
		{s + "some" + s + "random" + s + "file", s + "some" + s + "random" + s + "file"},
	}
	if runtime.GOOS == "linux" {
		current, _ := user.Current()
		ts = append(ts, []*struct {
			in  string
			out string
		}{
			{"~" + current.Username, HomeDir(current.Username)},
			{"~" + current.Username + s, HomeDir(current.Username) + s},
			{"~" + current.Username + s + "file", HomeDir(current.Username) + s + "file"},
		}...)
	}
	for _, c := range ts {
		out := ExpandTilde(c.in)
		t.Logf("%q => %q", c.in, out)
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("ExpandTilde(%q) = %q, want %q", c.in, out, c.out)
		}
	}
}
