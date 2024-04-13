BIN="bin"

default: all

.PHONY: all
all: fmt test build

.PHONY: build
build:
	go build -o $(BIN)/ ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race ./...

.PHONY: clean
clean:
	rm -rf $(BIN)
