all: build
  
clean:
	rm -f client

build: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

.PHONY: clean build
