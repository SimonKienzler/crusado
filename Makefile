TAG := $(shell git describe --tag --always)

.PHONY: build
build:
	@go build -o bin/crusado -ldflags="-X 'github.com/simonkienzler/crusado/cmd/version.CrusadoVersion=${TAG}'"
