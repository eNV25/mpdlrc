
export GOENV := "."

go        := "go"
gofmt     := "gofmt -r '(x) -> x' -s"
goimports := "goimports -local " + go-mod
gofumpt   := "gofumpt"

go-mod   := `go list -m`
go-files := "'" + replace(replace(`go list -f '{{ $d := .Dir }}{{ range .GoFiles }}{{ printf "%s/%s\n" $d . }}{{ end }}' ./...`, "'", "'\\''"), "\n", "' '") + "'"

build:
	{{ go }} build -v -o ./bin/ ./ ./cmd/...

run *args:
	{{ go }} run -v . {{ args }}

debug *args:
	{{ go }} run -v -tags=debug . {{ args }}

objdump bin sym='""':
	{{ go }} tool objdump -s {{ sym }} -S bin/{{ bin }} | less

test:
	{{ go }} test -v ./...

generate: build && fmt
	{{ go }} generate -v ./...

fmt:
	{{ go }} mod tidy
	{{ go }} fix ./...
	{{ go }} fmt ./...
	@ {{ gofmt }} -w -l {{ go-files }}
	@ {{ goimports }} -w -l {{ go-files }}
	@ {{ gofumpt }} -w -l {{ go-files }}

checkfmt:
	@ ! [ "$({{ gofmt }} -l {{ go-files }} | wc -l)" -gt 0 ]
	@ ! [ "$({{ goimports }} -l {{ go-files }} | wc -l)" -gt 0 ]
	@ ! [ "$({{ gofumpt }} -l {{ go-files }} | wc -l)" -gt 0 ]

install-tools should="! command -v":
	{{ should }} goimports && go install -v golang.org/x/tools/cmd/goimports@latest || true
	{{ should }} gofumpt && go install -v mvdan.cc/gofumpt@latest || true
