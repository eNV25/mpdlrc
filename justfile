
go        := "go"
gofmt     := "gofmt -r '(x) -> x' -s"
goimports := "goimports -local " + go-mod
gofumpt   := "gofumpt"

go-mod   := `go list -m`
go-files := replace(```go list -f '{{- $d := .Dir -}}{{- range .GoFiles -}}{{- printf "%s/%s" $d . | printf "%q\n" -}}{{- end -}}' ./...```, "\n", " ")

build:
	{{ go }} build -v -o ./bin/ ./...

run:
	{{ go }} run -v .

debug:
	{{ go }} run -v -tags=debug .

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

checkfmt:
	@ ! [ "$({{ gofmt }} -l {{ go-files }} | wc -l)" -gt 0 ]
	@ ! [ "$({{ goimports }} -l {{ go-files }} | wc -l)" -gt 0 ]
	@ ! [ "$({{ gofumpt }} -l {{ go-files }} | wc -l)" -gt 0 ]

install-tools FORCE="! command -v goimports":
	{{ FORCE }} && go install -v golang.org/x/tools/cmd/goimports@latest || true
	{{ FORCE }} && go install -v mvdan.cc/gofumpt@latest || true

