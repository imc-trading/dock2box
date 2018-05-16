SETUPIMG=alpine
REL=$(shell git rev-parse --verify --short HEAD)
CODENAME=Dapper Snapper

all: docker-build

clean:
	rm -f initrd kernel redhat-release agent
#	rm -rf openonload-*
	rm -rf .certs-cache

docker-clean:
	docker rmi ${IMG}:dev ${IMG}:prod &>/dev/null || true
	docker rm -f setup &>/dev/null || true

docker-build:
	cd ../client && make build
	cp ../client/client ./client
	rsync -a ../certs/ .certs-cache
	docker pull $(SETUPIMG) &>/dev/null
	docker run --rm --name=setup -v "$${PWD}":/build -w /build $(SETUPIMG) ./setup.sh $(REL) "$(CODENAME)"
	docker build --pull=true --no-cache .

.PHONY: clean docker-clean docker-build
