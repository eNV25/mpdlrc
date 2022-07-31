
go.module != go list -m

build: .phony
	go build -v -o ./bin/ ./...

run: .phony
	go run -v .

debug: .phony
	go run -v -tags=debug .

test: .phony
	go test -v ./...

gen: generate fmt .phony

generate: .phony
	go generate -v ./...

fmt: .phony
	go mod tidy
	go fix ./...
	go fmt ./...
	gofmt -r '(x) -> x' -s -w -l .
	goimports -local '${go.module}' -w -l .
	gofumpt -w -l .

checkfmt: .phony
	! [ "$$(gofmt -r '(x) -> x' -s -l . | wc -l)" -gt 0 ]
	! [ "$$(goimports -local '${go.module}' -l . | wc -l)" -gt 0 ]
	! [ "$$(gofumpt -l . | wc -l)" -gt 0 ]

.PHONY: .phony
.phony:
