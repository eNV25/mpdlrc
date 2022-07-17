//go:build generate
// +build generate

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	MPD_BUNDLE_PKG = "github.com/env25/gompd/v2/mpd@my"
	MPD_BUNDLE_OUT = "mpd_bundle.go"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func run(capture bool, command string, args ...string) (out []byte) {
	println("+", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	var err error
	if capture {
		out, err = cmd.Output()
	} else {
		cmd.Stdout = os.Stderr
		err = cmd.Run()
	}
	check(err)
	return
}

func main() {
	pkgname, err := exec.Command("go", "list").Output()
	check(err)
	pkgbase := filepath.Base(strings.TrimSpace(string(pkgname)))

	// initialise this file, otherwise bundle errors out
	check(os.WriteFile(MPD_BUNDLE_OUT, []byte("package "+pkgbase+"\n"), 0o644))

	// get deps
	run(false, "go", "get", "-v", MPD_BUNDLE_PKG)

	// run bundle without -o so that is doesn't add a go:generate line
	check(os.WriteFile(MPD_BUNDLE_OUT,
		run(true, "bundle", strings.Split(MPD_BUNDLE_PKG, "@")[0]), 0o644))

	// clean deps
	run(false, "go", "mod", "tidy", "-v")
}
