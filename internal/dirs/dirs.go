package dirs

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// GetEnv is a wrapper of [os.Getenv] implementing fallback for some environment variables.
func GetEnv(key string) string {
	switch key {
	case "HOME":
		return HomeDir("")
	}
	ret := os.Getenv(key)
	if strings.HasPrefix(key, "XDG_") {
		if ret != "" {
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
	}
	return ret
}

// ExpandEnv works like [os.ExpandEnv] but uses [GetEnv].
func ExpandEnv(s string) string {
	return os.Expand(s, GetEnv)
}

// RootDir returns the system root directory.
func RootDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("SYSTEMDRIVE") + string(os.PathSeparator)
	}
	return string(os.PathSeparator)
}

// HomeDir returns the user's home directory.
func HomeDir(usr string) string {
	var u *user.User
	var err error
	if usr == "" {
		if h, err := os.UserHomeDir(); err == nil {
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

// ExpandTilde expands tilde "~" into the home directory.
func ExpandTilde(s string) string {
	const pathSeparator = string(os.PathSeparator)
	s = filepath.FromSlash(s)
	if strings.HasPrefix(s, "~") {
		u, p, sep := strings.Cut(s[1:], pathSeparator)
		if sep {
			// ~/path or ~user/path
			return HomeDir(u) + pathSeparator + p
		}
		// ~ or ~user
		return HomeDir(u)
	}
	// path or /path
	return s
}
