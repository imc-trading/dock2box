NAME=alpine3.4.3
REPO=dock2box
REG=
REL=$(shell git rev-parse --verify --short HEAD)

all: build

clean:
	rm -f initrd kernel release

build:
	echo $(REL) >release
	docker pull $(REG)$(REPO)/$(NAME):latest 
	docker run --rm -v "$${PWD}":/build -w /build $(REG)$(REPO)/$(NAME):latest ./setup.sh
