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
		{"/", "/"},
		{"~", HomeDir("")},
		{"~/", HomeDir("") + "/"},
		{"~/file", HomeDir("") + "/file"},
		{"~/directory/", HomeDir("") + "/directory/"},
		{"~/directory/file", HomeDir("") + "/directory/file"},
		{"~" + current.Username, HomeDir(current.Username)},
		{"~" + current.Username + "/", HomeDir(current.Username) + "/"},
		{"~" + current.Username + "/file", HomeDir(current.Username) + "/file"},
		{"~" + current.Username + "/directory/", HomeDir(current.Username) + "/directory/"},
		{"file", "file"},
		{"some/random/file", "some/random/file"},
		{"/some/random/file", "/some/random/file"},
	} {
		out := ExpandTilde(c.in)
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("ExpandTilde(%q) = %q, want %q", c.in, out, c.out)
		}
	}
}
