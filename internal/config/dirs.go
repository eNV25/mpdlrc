package config

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/env25/mpdlrc/internal/ustrings"
)

func ConfigDir(usr string) string {
	if c, ok := os.LookupEnv("XDG_CONFIG_HOME"); !ok {
		return filepath.Join(HomeDir(usr), ".config")
	} else {
		return c
	}
}

func HomeDir(usr string) (h string) {
	var u *user.User
	var err error
	if usr == "" {
		h, err = os.UserHomeDir()
		if err == nil {
			return
		}
		u, err = user.Current()
	} else {
		u, err = user.Lookup(usr)
	}
	if err == nil {
		h = u.HomeDir
		return
	}
	return
}

func ExpandTilde(str string) string {
	switch {
	case str == "":
		return ""
	case strings.HasPrefix(str, "~"):
		// ~ or ~/path or ~user/path
		u, p, _ := ustrings.Cut(str[1:], string(filepath.Separator))
		return filepath.Join(HomeDir(u), p) // calls filepath.Clean
	default:
		// path or /path
		return filepath.Clean(str)
	}
}
