VERSION  := 1.22
BIN_DIR  := bin
BIN_NAME := ranga-$(VERSION)

GCFLAGS  := -gcflags="all=-B"
LDFLAGS  := -ldflags="-s -w -X main.version=$(VERSION)"

.DEFAULT_GOAL := help

.PHONY: help build run bench test cover race clean \
        build-docker build-linux build-windows build-mac

help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build:
	go build $(GCFLAGS) $(LDFLAGS) -o $(BIN_DIR)/$(BIN_NAME) ./cmd
	chmod +x $(BIN_DIR)/$(BIN_NAME)

run: build
	./$(BIN_DIR)/$(BIN_NAME)

bench: build
	./$(BIN_DIR)/$(BIN_NAME) bench

test:
	go test ./...

cover:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out

race:
	go test --race ./...

clean:
	rm -rf $(BIN_DIR)/$(BIN_NAME)
	go clean -cache -testcache -modcache

build-docker:
	docker build --build-arg VERSION=$(VERSION) --target exporter --output ./$(BIN_DIR) .

define cross_build
	env CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build $(GCFLAGS) -ldflags="-s -w -X 'main.version=$(VERSION)'" -o $(BIN_DIR)/$(BIN_NAME)-$(1)-$(2)$(3) ./cmd
endef

build-linux:
	$(call cross_build,linux,amd64,)

build-windows:
	$(call cross_build,windows,amd64,.exe)

build-mac:
	$(call cross_build,darwin,arm64,)
