
build:
	go build -v $(GOFLAGS) -o ./bin/ ./...

test:
	go test -v $(GOFLAGS) ./...

fmt:
	goimports -local github.com/env25/mpdlrc -w -l .
	gofumpt -w -l .

