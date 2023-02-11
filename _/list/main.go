package main

import (
	"flag"
	"fmt"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	flag_m       = flag.Bool("m", false, "print module names")
	flag_gofiles = flag.Bool("gofiles", false, "print go files")
	flag_sh      = flag.Bool("sh", false, "print outputs space separated and quoted in Bourne shell syntax")
	flag_rc      = flag.Bool("rc", false, "print outputs space separated and quoted in rc (or powershell) shell syntax")
)

func main() {
	flag.Parse()

	var cfg packages.Config
	cfg.Mode = packages.NeedFiles | packages.NeedModule

	pkgs, err := packages.Load(&cfg, flag.Args()...)
	check(err)

	var outputs []string

	switch {
	case *flag_m:
		for _, pkg := range pkgs {
			outputs = append(outputs, pkg.Module.Path)
		}
	case *flag_gofiles:
		for _, pkg := range pkgs {
			outputs = append(outputs, pkg.GoFiles...)
		}
	default:
		flag.Usage()
		return
	}

	for i := range outputs {
		outputs[i] = quote(outputs[i])
	}

	if *flag_rc || *flag_sh {
		fmt.Println(strings.Join(outputs, " "))
	} else {
		for i := range outputs {
			fmt.Println(outputs[i])
		}
	}
}

func quote(s string) string {
	switch {
	case *flag_sh:
		return `'` + strings.ReplaceAll(s, `'`, `'\''`) + `'`
	case *flag_rc:
		return `'` + strings.ReplaceAll(s, `'`, `''`) + `'`
	}
	return s
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
