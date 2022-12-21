//go:build debug

package config

import (
	"log"
	"net/http"
	_ "net/http/pprof" // enable pprof for debugging
)

// Debug is true for debug builds.
const Debug = true

func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
