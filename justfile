
export GOENV := "."

go        := "go"
gofmt     := "gofmt -r '(x) -> x' -s"
goimports := "go run golang.org/x/tools/cmd/goimports -local " + go-mod
gofumpt   := "go run mvdan.cc/gofumpt"

go-mod   := `go run _/list -sh -m`
go-files := `go run _/list -sh -gofiles ./... ./tools/...`

_ := `go mod tidy && cd _ && go mod tidy`

build:
	{{ go }} build -v -o ./bin/ .

run *args:
	{{ go }} run -v . {{ args }}

debug *args:
	{{ go }} run -v -tags=debug . {{ args }}

objdump bin sym='""':
	{{ go }} tool objdump -s {{ sym }} -S bin/{{ bin }} | less

list-inline:
	go build -gcflags='-m' -o /dev/null ./... 2>&1 | grep -e 'can inline' -e 'inlining'

test:
	{{ go }} test -v ./...

generate: && fmt
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
