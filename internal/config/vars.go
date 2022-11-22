package config

import (
	"path/filepath"

	"github.com/env25/mpdlrc/internal/dirs"
)

// DefaultFiles returns the default configuration files.
func DefaultFiles() []string {
	return []string{
		filepath.Join(dirs.GetEnv("XDG_CONFIG_HOME"), "mpdlrc", "config.toml"),
	}
}
