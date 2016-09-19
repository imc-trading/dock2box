NAME=alpine3.4
ORG=dock2box
REL=$(shell git rev-parse --verify --short HEAD)

all: build

clean:
	rm -f initrd kernel release

build:
	echo $(REL) >release
	docker pull $(ORG)/$(NAME):latest
	docker run --rm -v "$${PWD}":/build -w /build $(ORG)/$(NAME):latest ./setup.sh
