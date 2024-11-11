GO_FLAGS = -ldflags="-s -w"
BUILD_OUT = ./dist/thunder

.PHONY: all build build-front build-golang compress-binary clean dev

all: build dev

build: build-front build-golang compress-binary

dev:
	cd app && go run -tags debug .

build-front:
	 cd ./modules/frontend &&  yarn && yarn generate

build-golang:
	cd app && CGO_ENABLED=0 go build $(GO_FLAGS) -o "../$(BUILD_OUT)"

compress-binary:
	upx -9 -k $(BUILD_OUT)

clean:
	rm -f $(BUILD_OUT)