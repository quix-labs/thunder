GO_FLAGS = -ldflags="-s -w"
BUILD_OUT = ./dist/thunder

.PHONY: all build build-front build-golang clean

all: build

build: build-front build-golang


build-front:
	 cd ./modules/frontend && yarn  && yarn generate

build-golang:
	cd app && CGO_ENABLED=0 go build $(GO_FLAGS) -o "../$(BUILD_OUT)"

clean:
	rm -f $(BUILD_OUT)