
export GOENV := "."

go        := "go"
gofmt     := "gofmt -r '(x) -> x' -s"
goimports := "goimports -local " + go-mod
gofumpt   := "gofumpt"

go-mod   := `go list -m`
go-files := "'" + replace(replace(`go list -f '{{ $d := .Dir }}{{ range .GoFiles }}{{ printf "%s/%s\n" $d . }}{{ end }}' ./...`, "'", "'\\''"), "\n", "' '") + "'"

build:
	{{ go }} build -v -o ./bin/ ./...

run:
	{{ go }} run -v .

debug *args:
	{{ go }} run -v -tags=debug . {{ args }}

objdump bin sym='""':
	{{ go }} tool objdump -s {{ sym }} -S bin/{{ bin }} | less

test:
	{{ go }} test -v ./...

generate: install-tools build && fmt
	{{ go }} generate -v ./...

fmt:
	{{ go }} mod tidy
	{{ go }} fix ./...
	{{ go }} fmt ./...
	@ {{ gofmt }} -w -l {{ go-files }}
	@ {{ goimports }} -w -l {{ go-files }}
	@ {{ gofumpt }} -w -l {{ go-files }}

checkfmt: install-tools
	@ ! [ "$({{ gofmt }} -l {{ go-files }} | wc -l)" -gt 0 ]
	@ ! [ "$({{ goimports }} -l {{ go-files }} | wc -l)" -gt 0 ]
	@ ! [ "$({{ gofumpt }} -l {{ go-files }} | wc -l)" -gt 0 ]

install-tools should="! command -v goimports":
	{{ should }} && go install -v golang.org/x/tools/cmd/goimports@latest || true
	{{ should }} && go install -v mvdan.cc/gofumpt@latest || true

