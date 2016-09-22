NAME=alpine3.4.3
REPO=dock2box
REL=$(shell git rev-parse --verify --short HEAD)

all: publish

clean:
	rm -f initrd kernel release

build:
	echo $(REL) >release
	docker pull $(REPO)/$(NAME):latest 
	docker run --rm -v "$${PWD}":/build -w /build $(REPO)/$(NAME):latest ./setup.sh
