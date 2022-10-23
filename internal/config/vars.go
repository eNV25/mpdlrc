package config

import (
	"path/filepath"
)

// DefaultFiles returns the default configuration files.
func DefaultFiles() []string {
	return []string{
		filepath.Join(GetEnv("XDG_CONFIG_HOME"), "mpdlrc", "config.toml"),
	}
}
