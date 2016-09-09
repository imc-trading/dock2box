NAME=alpine3.4
ORG=dock2box

all: build

clean:
	rm -f initrd kernel apkovl.tar.gz

build:
	docker pull $(USER)/$(NAME):latest
	docker run --rm -v "$${PWD}":/build -w /build $(ORG)/$(NAME):latest ./setup.sh
