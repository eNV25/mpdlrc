
go.module     != go list -m
cmd.go        := go
cmd.gofmt     := gofmt -r '(x) -> x' -s
cmd.goimports := goimports -local '${go.module}'
cmd.gofumpt   := gofumpt

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
	${cmd.gofmt} -w -l .
	${cmd.goimports} -w -l .
	${cmd.gofumpt} -w -l .

checkfmt: .phony
	! [ "$$(${cmd.gofmt} -l . | wc -l)" -gt 0 ]
	! [ "$$(${cmd.goimports} -l . | wc -l)" -gt 0 ]
	! [ "$$(${cmd.gofumpt} -l . | wc -l)" -gt 0 ]

gen: build generate fmt .phony

tools:
	go install -v golang.org/x/tools/cmd/goimports@latest
	go install -v mvdan.cc/gofumpt@latest

.PHONY: .phony
.phony:
