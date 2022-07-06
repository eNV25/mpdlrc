//go:build debug
// +build debug

package config

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

const Debug = true

func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
