GO_FLAGS = -ldflags="-s -w"
BUILD_OUT = ./dist/thunder
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: all build build-golang compress-binary clean dev dev-nodebug

all: clean build

build: build-golang compress-binary

dev:
	cd dev && go run -tags debug .

dev-nodebug:
	cd dev && go run .

generate:
	go generate ./...

build-golang:
	cd app && CGO_ENABLED=0 go build $(GO_FLAGS) -o "../$(BUILD_OUT)"

compress-binary:
	upx --best --lzma $(BUILD_OUT)

clean:
	rm -rf $(BUILD_OUT)