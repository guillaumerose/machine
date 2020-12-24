.PHONY: all
all: vendor lint test

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test ./...
