
all: lint build

.PHONY: build
build:
	@go build -o ./ ./...

.PHONY: lint
lint:
	golangci-lint run
