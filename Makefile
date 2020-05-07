all: fmt build

build: clean fmt
	CGO_ENABLED=0 go build -o bin/webserver -a

clean:
	rm -rf bin
	
fmt:
	go fmt ./...

.PHONY: all build clean fmt 
