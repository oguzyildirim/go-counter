BUILD_DIR := $(shell pwd)/_out

.PHONY: build
all: build

build:
	cd cmd/server && go build -o $(BUILD_DIR)/server

remove:
	rm -rf $(BUILD_DIR)

test:
	go test -race ./...
