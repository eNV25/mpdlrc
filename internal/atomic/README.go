//go:build generate
// +build generate

// TODO:
// This was forked from go.uber.org/atomic in https://github.com/eNV25/atomic.
// We use this since changes weren't merged
// and we don't want to use mod replace.

package main

import (
	"os"
	"path/filepath"
	"regexp"
)

//go:generate git clone -b my --depth 1 --single-branch https://github.com/eNV25/atomic.git
//go:generate go run README.go

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	files := make(map[string]struct{})

	mtchs, err := filepath.Glob(filepath.FromSlash("atomic/*.go"))
	check(err)
	for _, f := range mtchs {
		files[f] = struct{}{}
	}

	mtchs, err = filepath.Glob(filepath.FromSlash("atomic/*.s"))
	check(err)
	for _, f := range mtchs {
		files[f] = struct{}{}
	}

	mtchs, err = filepath.Glob(filepath.FromSlash("atomic/*_test.go"))
	check(err)
	for _, f := range mtchs {
		delete(files, f)
	}

	mtchs, err = filepath.Glob(filepath.FromSlash("atomic/gen*.go"))
	check(err)
	for _, f := range mtchs {
		delete(files, f)
	}

	go_gen_line := regexp.MustCompile("(?m)[\r\n]+^//go:generate.*$")
	for f := range files {
		dst := filepath.Base(f)
		data, err := os.ReadFile(f)
		check(err)
		data = go_gen_line.ReplaceAll(data, nil) // delete
		err = os.WriteFile(dst, data, 0o644)
		check(err)
	}

	err = os.RemoveAll(filepath.FromSlash("atomic/"))
	check(err)
}
