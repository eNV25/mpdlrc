package config

import (
	"os"
	"os/user"
	"path"
)

func ConfigDir() string {
	if c, ok := os.LookupEnv("XDG_CONFIG_HOME"); !ok {
		return path.Join(HomeDir(), ".config")
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
	if u, err := user.Lookup(usr); err != nil {
		// fallback
		// path.Dir("/home/user") => "/home"
		return path.Join(path.Dir(HomeDir()), usr)
	} else {
		return u.HomeDir
	}
}
