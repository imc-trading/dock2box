ifeq ($(REL), 7.0)
IMAGE=centos7.0.1406
endif

ifeq ($(REL), 7.1)
IMAGE=centos7.1.1503
endif

ifeq ($(REL), 7.2)
IMAGE=centos7.2.1511
endif

ifndef REL
REL=7
IMAGE=centos7
endif

ifdef NOCACHE
OPTS=--no-cache
endif

USER=dock2box
NAME=centos-${REL}

all: build

clean:
	docker rmi ${NAME} &>/dev/null || true

build: clean
	sed 's/centos:centos7/centos:${IMAGE}/' Dockerfile >.Dockerfile.tmp
	docker build --pull=true ${OPTS} -t ${USER}/${NAME}:latest -f .Dockerfile.tmp .

push: build
	docker login
	docker push ${USER}/${NAME}:latest
