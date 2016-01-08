USER=imc-trading
NAME=d2bsrv
SRC=src/github.com/${USER}/${NAME}
TAG:=$(shell awk -F '"' '/Version/ {print $$2}' ${SRC}/version.go)
BUILD=.build

all: build

clean:
	rm -rf pkg bin ${BUILD}

build:
	gb build all

docker-clean:
	docker rmi ${NAME} &>/dev/null || true

docker-build: docker-clean
	docker run --rm -it -v "$$PWD":/go -w /go mickep76/alpine-golang:latest
	docker build --pull=true --no-cache -t ${USER}/${NAME}:${TAG} .

docker-push: docker-build
	docker login -u ${USER}
	docker push ${USER}/${NAME}:${TAG}

docker-push-latest: docker-push
	docker tag -f ${USER}/${NAME}:${TAG} ${USER}/${NAME}:latest
	docker push ${USER}/${NAME}:latest
