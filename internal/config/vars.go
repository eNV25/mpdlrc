package config

import (
	"path/filepath"
)

func ConfigFiles() []string {
	return []string{
		filepath.Join(GetEnv("XDG_CONFIG_HOME"), "mpdlrc", "config.toml"),
	}
}
