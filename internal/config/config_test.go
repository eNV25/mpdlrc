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
		{"~", HomeDir()},
		{"~/", HomeDir()},
		{"~/directory/", HomeDir() + "/directory"},
		{"~/directory/file", HomeDir() + "/directory/file"},
		{"~" + current.Username + "/directory", HomeDirUser(current.Username) + "/directory"},
	} {
		out := expandTilde(c.in)
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("expandTilde(%q) => %q, expected => %q", c.in, out, c.out)
		}
	}
}
