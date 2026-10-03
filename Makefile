EXE     ?= ranga
VERSION ?= 1.14

all:
	go build -ldflags="-X main.version=$(VERSION)" -o $(EXE) ./cmd
