package config

import (
	"os"
	"os/user"
	"path/filepath"
)

func ConfigDir() string {
	if c, ok := os.LookupEnv("XDG_CONFIG_HOME"); !ok {
		return filepath.Join(HomeDir(), ".config")
	} else {
		return c
	}
}

func HomeDir() string {
	if h, err := os.UserHomeDir(); err != nil {
		if u, e := user.Current(); e != nil {
			panic(err)
		} else {
			return u.HomeDir
		}
	} else {
		return h
	}
}

func HomeDirUser(usr string) string {
	if usr == "" {
		return HomeDir()
	}
	if u, err := user.Lookup(usr); err != nil {
		// fallback
		// path.Dir("/home/user") => "/home"
		return filepath.Join(filepath.Dir(HomeDir()), usr)
	} else {
		return u.HomeDir
	}
}
