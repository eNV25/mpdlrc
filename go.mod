module github.com/env25/mpdlrc

go 1.16

// replace version with `master`
// then run `go mod tidy`
require github.com/gdamore/tcell/v2 v2.3.12-0.20210612024312-b60a903b9868

require (
	github.com/fhs/gompd/v2 v2.2.0
	github.com/neeharvi/textwidth v0.0.0-20210704185812-614098cf670d
	github.com/pelletier/go-toml/v2 v2.0.0-beta.3
	github.com/spf13/pflag v1.0.5
	golang.org/x/text v0.3.6
)
