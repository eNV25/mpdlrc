<!-- vi: set wrap linebreak: -->

# mpdlrc

https://github.com/eNV25/mpdlrc/

`mpdlrc` displays synchronized lyrics for the currently playing track in the Music Player Daemon (MPD). It uses the track's file path to find an `.lrc` file, e.g. `file.mp3 => file.lrc`. In the future, it may be extended to support synchronized lyrics embedded in the audio file.

## Install

Installation requires a [go](http://golang.org/) compiler and the `go` tool.

Install to `${GOPATH:-$HOME/go}/bin` using the `go` tool.

```console
$ go install github.com/env25/mpdlrc@latest
```

or use master branch

```console
$ go install github.com/env25/mpdlrc@master
```

NOTE: You may need to add `${GOPATH:-$HOME/go}/bin` to `$PATH`.

## Set up and Configure

If you run mpd in your local machine and you use `MPD_HOST` (and `MPD_PORT`) (see [man:mpc(1)]), you need no further configuration. If connecting using a UNIX socket, `mpdlrc` queries for the music directory using the MPD protocol. Otherwise, because of a restriction in the MPD protocol, it reads `mpd` configuration files as fallback. `mpdlrc` assumes your lyrics are stored alongside your music with the same filepath minus extension.

If you run `mpd` on a remote machine, you should explicitly configure the lyrics directory with `mpdlrc`. Currently you require a directory path accessible using the local machine, so you will probably need to mount your files using something like `sshfs` or `rclone`. If using a different directory, the layout should match what is used by `mpd`.

Any automatic configuration can be overridden using a TOML config file. The config file should be located in `${XDG_CONFIG_HOME:-$HOME/.config}/mpdlrc/config.toml`. More exhaustive documentation for the config file can be found in [docs/config-docs.toml](docs/config-docs.toml).

## Run

```console
$ mpdlrc
```

## Screenshot

![screenshot.png](https://user-images.githubusercontent.com/61089994/178155519-89f2829c-9640-459b-8df0-1478354e26ab.png)

[man:mpc(1)]: https://man.archlinux.org/man/mpc.1
