BINARY_NAME=goseal

build:
	CGO_ENABLED=0 go build -ldflags '-w -s' -trimpath -o $(BINARY_NAME) cmd/goseal/main.go

test:
	go test -race ./...

test-coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic -coverpkg=./... ./...

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

fmt:
	golangci-lint fmt

.PHONY: build test test-coverage lint lint-fix fmt
