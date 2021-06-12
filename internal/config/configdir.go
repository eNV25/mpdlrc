package config

import (
	"os"
	"path"
)

func ConfigDir() string {
	if c, ok := os.LookupEnv("XDG_CONFIG_HOME"); !ok {
		return path.Join(HomeDir(), ".config")
	} else {
		return c
	}
}
