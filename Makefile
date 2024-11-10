GO_FLAGS = -ldflags="-s -w"
BUILD_OUT = ./dist/thunder

.PHONY: all build build-front build-golang compress-binary clean

all: build

build: build-front build-golang build-elastic compress-binary


build-front:
	 cd ./modules/frontend &&  yarn && yarn generate

build-golang:
	cd app && CGO_ENABLED=0 go build $(GO_FLAGS) -o "../$(BUILD_OUT)"

build-elastic:
	 cd target-drivers/elastic &&  go build $(GO_FLAGS) -buildmode=plugin -o elastic.so .

compress-binary:
	upx -9 -k $(BUILD_OUT)

clean:
	rm -f $(BUILD_OUT)