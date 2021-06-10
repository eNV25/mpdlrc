package config

import "path"

var (
	ConfigFiles = []string{
		path.Join(ConfigDir(), "mpdlrc/config.toml"),
	}
)
