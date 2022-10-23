
go.module     != go list -m
cmd.go        := go
cmd.gofmt     := gofmt -r '(x) -> x' -s
cmd.goimports := goimports -local '${go.module}'
cmd.gofumpt   := gofumpt

go.files != go list -f '{{ $$d := .Dir }}{{ range .GoFiles }}{{ printf "%s/%s" $$d . | printf "%q\n" }}{{ end }}' ./...

build: .phony
	${cmd.go} build -v -o ./bin/ ./...

run: .phony
	${cmd.go} run -v .

debug: .phony
	${cmd.go} run -v -tags=debug .

test: .phony
	${cmd.go} test -v ./...

generate: .phony
	${cmd.go} generate -v ./...

fmt: .phony
	${cmd.go} mod tidy
	${cmd.go} fix ./...
	${cmd.go} fmt ./...
	@ ${cmd.gofmt} -w -l ${go.files}
	@ ${cmd.goimports} -w -l ${go.files}
	@ ${cmd.gofumpt} -w -l ${go.files}

checkfmt: .phony
	@ ! [ "$$(${cmd.gofmt} -l ${go.files} | wc -l)" -gt 0 ]
	@ ! [ "$$(${cmd.goimports} -l ${go.files} | wc -l)" -gt 0 ]
	@ ! [ "$$(${cmd.gofumpt} -l ${go.files} | wc -l)" -gt 0 ]

gen: build generate fmt .phony

tools:
	go install -v golang.org/x/tools/cmd/goimports@latest
	go install -v mvdan.cc/gofumpt@latest

.PHONY: .phony
.phony:
