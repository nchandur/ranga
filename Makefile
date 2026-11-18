VERSION  := 1.22
BIN_DIR  := bin
BIN_NAME := ranga-$(VERSION)

EXE      ?= $(BIN_DIR)/$(BIN_NAME)

GCFLAGS  := -gcflags="all=-B"
LDFLAGS  := -ldflags="-s -w -X main.version=$(VERSION)"

.DEFAULT_GOAL := build

.PHONY: help build run bench test cover race clean docker linux windows mac

help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build:
	@mkdir -p $(dir $(EXE))
	CGO_ENABLED=0 go build $(GCFLAGS) $(LDFLAGS) -o $(EXE) ./cmd
	chmod +x $(EXE)

run: build
	./$(EXE)

bench: build
	./$(EXE) bench

test:
	go test ./...

cover:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out

race:
	go test --race ./...

clean:
	rm -rf $(BIN_DIR)/$(BIN_NAME)
	go clean -cache -testcache

docker:
	docker build --build-arg VERSION=$(VERSION) --target exporter --output ./$(BIN_DIR) .

define cross_build
	env CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build $(GCFLAGS) -ldflags="-s -w -X 'main.version=$(VERSION)'" -o $(BIN_DIR)/$(BIN_NAME)-$(1)-$(2)$(3) ./cmd
endef

linux:
	$(call cross_build,linux,amd64,)

windows:
	$(call cross_build,windows,amd64,.exe)

mac:
	$(call cross_build,darwin,arm64,)
