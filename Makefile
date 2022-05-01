
build:
	go build -v $(GOFLAGS) -o ./bin/ ./...

fmt:
	goimports -w -l .
	gofumpt -w -l .

