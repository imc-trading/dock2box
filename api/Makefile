NAME=etcdrest-config
VERSION=0.1
USER=imctrading
BUILD=.build

all: build

clean:
	rm -rf pkg bin ${BUILD}

build:
	gb build all

docker-clean: clean
	docker rmi ${NAME} &>/dev/null || true

docker:
	docker pull mickep76/alpine-golang:latest
	docker run --rm -it -v "$$PWD":/app -w /app mickep76/alpine-golang:latest
	docker build --pull=true --no-cache -t ${USER}/${NAME}:${VERSION} .
	docker tag -f ${USER}/${NAME}:${VERSION} ${USER}/${NAME}:latest

push:
	docker login -u ${USER}
	docker push ${USER}/${NAME}:${VERSION}
	docker push ${USER}/${NAME}:latest
