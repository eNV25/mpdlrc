# mpdlrc

https://github.com/eNV25/mpdlrc/

`mpdlrc` displays synchronized lyrics for the currently playing track. It uses the track's file path to find an `.lrc` file, e.g. `file.mp3 => file.lrc`. In the future, it may be extended to support synchronized lyrics embedded in the audio file.

## Install

Installation requires a [go](http://golang.org/) compiler and the `go` tool.

Install to `${GOPATH:-$HOME/go}/bin` using the `go` tool.
```
$ go install github.com/env25/mpdlrc@latest
```
or use master branch
```
$ go install github.com/env25/mpdlrc@master
```
NOTE: You may need to add `${GOPATH:-$HOME/go}/bin` to `$PATH`.

## Set up and Configure

You must set up MPD first. If you use `MPD_HOST` (and `MPD_PORT`) for
the mpc command-line client see [man:mpc(1)](https://man.archlinux.org/man/mpc.1),
mpdlrc will pick those up. If you use a unix socket to connect to mpd
no further configuration is required. Otherwise, since mpd doesn't allow
clients to query the information, you need to at least configure
the `MusicDir` option.

Configuration is done using a TOML config file. The config file should be
located in `${XDG_CONFIG_HOME:-$HOME/.config}/mpdlrc/config.toml`. More
exhaustive documentation for the config file can be found in
[docs/config-docs.toml](docs/config-docs.toml).

Example file, after setting `MPD_HOST=${XDG_RUNTIME_DIR}/mpd/socket`:

```toml
MusicDir = "$HOME/Music"
```

## Run

```
$ mpdlrc
```

## Screenshot

![screenshot.png](https://user-images.githubusercontent.com/61089994/178155519-89f2829c-9640-459b-8df0-1478354e26ab.png)
