.SILENT:

build: 
	go build -o bin/main

run: build
	./bin/main

test:
	go test

clean:
	go clean -cache