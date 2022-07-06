
go.module != go list -m

build:
	go build -v -o ./bin/ ./...

debug:
	go build -v -tags=debug -o ./bin/ ./...

run:
	go run -v .

run-debug:
	go run -v -tags=debug .

test:
	go test -v ./...

fmt:
	go mod tidy
	go fix ./...
	go fmt ./...
	gofmt -s -w -l .
	goimports -local '${go.module}' -w -l .
	gofumpt -w -l .

checkfmt:
	! [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]
	! [ "$$(goimports -local '${go.module}' -l . | wc -l)" -gt 0 ]
	! [ "$$(gofumpt -l . | wc -l)" -gt 0 ]

