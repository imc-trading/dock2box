all: test readme

format:
	gofmt -w=true .

syntax: format
	golint .
	go vet .

test: syntax
#	go test

readme:
	cp README_HEAD.md README.md
	godoc2md github.com/imc-trading/dock2box/client | grep -v Generated >>README.md
