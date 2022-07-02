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
	if strings.HasPrefix(str, "~") {
		u, p, sep := ustrings.Cut(str[1:], string(os.PathSeparator))
		if os.PathSeparator != '/' && !sep {
			u, p, sep = ustrings.Cut(str[1:], string('/'))
		}
		if sep {
			// ~/path or ~user/path
			return HomeDir(u) + string(os.PathSeparator) + p
		}
		// ~ or ~user
		return HomeDir(u)
	}
	// path or /path
	return str
}
