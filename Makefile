.PHONY: clean generate test build

BIN_OUTPUT := $(if $(filter $(shell go env GOOS), windows), dist/gci.exe, dist/gci)

default: clean generate test build

clean:
	@echo BIN_OUTPUT: ${BIN_OUTPUT}
	@rm -rf dist cover.out

VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
GIT_TAG := $(shell git describe --tags --abbrev=0 --exact-match > /dev/null 2>&1; echo $$?)

ifneq ($(GIT_TAG),0)
	ifeq ($(origin VERSION),file)
		VERSION := devel
	endif
endif

LDFLAGS ?= -w -X main.Version=${VERSION}
build: clean
	@go build -v -trimpath -ldflags "${LDFLAGS}" -o ${BIN_OUTPUT} .

test: clean
	@go test -v -count=1 -cover ./...

generate:
	@GOEXPERIMENT=arenas,boringcrypto,synctest go generate ./...
