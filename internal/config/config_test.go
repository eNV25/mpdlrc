package config

import (
	"reflect"
	"testing"
)

var expandTildeCases = []struct {
	in  string
	out string
}{
	{"~/directory", HomeDir() + "/directory"},
	{"~root/directory", HomeDirUser("root") + "/directory"},
}

func TestExpandTilde(t *testing.T) {
	for _, c := range expandTildeCases {
		exp := expandTilde(c.in)
		if !reflect.DeepEqual(exp, c.out) {
			t.Errorf("expandTilde(%q) => %q, expected => %q", c.in, exp, c.out)
		}
	}
}
