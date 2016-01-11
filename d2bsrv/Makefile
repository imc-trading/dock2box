USER=mickep76
NAME=d2bsrv
TAG=0.1
SRC=src/github.com/${USER}/dock2box/${NAME}
BUILD=.build

all: build

clean:
	rm -rf pkg bin ${BUILD}

build:
	go get github.com/constabulary/gb/...
	gb build all

docker-clean: clean
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
