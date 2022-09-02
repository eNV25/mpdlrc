module github.com/env25/mpdlrc

go 1.19

// use fork until PR merged
// https://github.com/fhs/gompd/pull/72
// https://github.com/eNV25/gompd/tree/my
require github.com/env25/gompd/v2 v2.2.1-0.20220711100057-a554ee3acd3d

require (
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/gdamore/tcell/v2 v2.5.2
	github.com/mattn/go-runewidth v0.0.13
	github.com/pelletier/go-toml/v2 v2.0.3
	github.com/rivo/uniseg v0.3.4
	go.uber.org/multierr v1.8.0
)

require (
	github.com/fhs/gompd/v2 v2.2.0 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20220318055525-2edf467146b5 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/text v0.3.7 // indirect
)
