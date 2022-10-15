package config

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func GetEnv(key string) string {
	switch key {
	case "HOME":
		return HomeDir("")
	}
	if strings.HasPrefix(key, "XDG_") {
		if ret := os.Getenv(key); ret != "" {
			return ret
		}
		switch key {
		case "XDG_CONFIG_HOME":
			return filepath.Join(HomeDir(""), ".config")
		case "XDG_CACHE_HOME":
			return filepath.Join(HomeDir(""), ".cache")
		case "XDG_DATA_HOME":
			return filepath.Join(HomeDir(""), ".local", "share")
		case "XDG_STATE_HOME":
			return filepath.Join(HomeDir(""), ".local", "state")
		case "XDG_DATA_DIRS":
			return filepath.Join(RootDir(), "usr", "local", "share") +
				":" + filepath.Join(RootDir(), "usr", "share")
		case "XDG_CONFIG_DIRS":
			return filepath.Join(RootDir(), "etc", "xdg")
		}
		return ""
	}
	return os.Getenv(key)
}

func ExpandEnv(s string) string {
	return os.Expand(s, GetEnv)
}

func RootDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("SYSTEMDRIVE") + string(os.PathSeparator)
	}
	return string(os.PathSeparator)
}

func HomeDir(usr string) string {
	var u *user.User
	var err error
	if usr == "" {
		h := ""
		h, err = os.UserHomeDir()
		if err == nil {
			return h
		}
		u, err = user.Current()
	} else {
		u, err = user.Lookup(usr)
	}
	if err == nil && u != nil {
		return u.HomeDir
	}
	return ""
}

func ExpandTilde(s string) string {
	s = filepath.FromSlash(s)
	if strings.HasPrefix(s, "~") {
		u, p, sep := strings.Cut(s[1:], string(os.PathSeparator))
		if sep {
			// ~/path or ~user/path
			return HomeDir(u) + string(os.PathSeparator) + p
		}
		// ~ or ~user
		return HomeDir(u)
	}
	// path or /path
	return s
}
