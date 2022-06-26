package config

import (
	"os/user"
	"reflect"
	"testing"
)

func TestExpandTilde(t *testing.T) {
	current, _ := user.Current()

	for _, c := range [...]*struct {
		in  string
		out string
	}{
		{"", ""},
		{"~", HomeDir("")},
		{"~/", HomeDir("")},
		{"~/directory/", HomeDir("") + "/directory"},
		{"~/directory/file", HomeDir("") + "/directory/file"},
		{"~" + current.Username + "/directory", HomeDir(current.Username) + "/directory"},
		{"/", "/"},
		{"///////", "/"},
		{"/some/random/dir/", "/some/random/dir"},
	} {
		out := ExpandTilde(c.in)
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("expandTilde(%q) = %q, want %q", c.in, out, c.out)
		}
	}
}
