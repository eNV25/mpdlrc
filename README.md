# mpdlrc

## Install

Installation requires a [go](http://golang.org/) compiler and the `go` tool.

Install to `${GOPATH:-$HOME/go}/bin` using the `go` tool.

    $ go install github.com/env25/mpdlrc@latest

or use master branch

    $ go install github.com/env25/mpdlrc@master

NOTE: You may need to add `${GOPATH:-$HOME/go}/bin` to `$PATH`.

## Setup and Configure

You'll need to setup MPD first.

Configuration is done using a TOML config file. The config file should be located in
`${XDG_CONFIG_HOME:-$HOME/.config}/mpdlrc/config.toml`.

Documentation for the config file can be found in [docs/config-docs.toml](docs/config-docs.toml).

Example file:

```toml
MusicDir = "$HOME/Music"
LyricsDir = "$HOME/Music"

[MPD]
Protocol = "unix"
Address = "${XDG_RUNTIME_DIR}/mpd/socket"
```

## Run

    $ mpdlrc

