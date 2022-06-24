package config

import (
	"os"
	"path/filepath"
)

func ConfigFiles() []string {
	return []string{
		filepath.Join(ConfigDir(""), "mpdlrc"+string(os.PathSeparator)+"config.toml"),
	}
}
