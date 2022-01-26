//go:build debug

package config

import (
	_ "net/http/pprof"
)

const Debug = true
