package config

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/env25/mpdlrc/internal/stringsu"
)

func ConfigDir() string {
	if c, ok := os.LookupEnv("XDG_CONFIG_HOME"); !ok {
		return filepath.Join(HomeDir(""), ".config")
	} else {
		return c
	}
}

func HomeDir(usr string) string {
	if usr == "" {
		h, err := os.UserHomeDir()
		if err != nil {
			u, errr := user.Current()
			if errr != nil {
				panic(err)
			}
			return u.HomeDir
		}
		return h
	}
	u, err := user.Lookup(usr)
	if err != nil {
		// fallback
		// return path.Dir("/home/current") + "/user"
		return filepath.Join(filepath.Dir(HomeDir("")), usr)
	}
	return u.HomeDir
}

func ExpandTilde(str string) string {
	switch {
	case strings.HasPrefix(str, "~"):
		// ~ or ~/path or ~user/path
		u, p, _ := stringsu.Cut(str[1:], string(os.PathSeparator))
		return filepath.Join(HomeDir(u), p) // calls filepath.Clean
	default:
		// path or /path
		return str
	}
}
