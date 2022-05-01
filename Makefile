
build:
	go build -v $(GOFLAGS) -o ./bin/ ./...

test:
	go test -v $(GOFLAGS) ./...

fmt:
	goimports -w -l .
	gofumpt -w -l .

