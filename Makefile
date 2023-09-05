TAG := $(shell git describe --tag --always)

.PHONY: build
build:
	@go build -o bin/crusado -ldflags="-X 'github.com/simonkienzler/crusado/cmd/version.CrusadoVersion=${TAG}'"

download: ## Downloads the dependencies
	@go mod download

lint: download ## Lints all code with golangci-lint
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run