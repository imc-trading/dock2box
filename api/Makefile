REPO=dock2box
NAME=api
TAG=0.1
SRC=src/github.com/${REPO}/dock2box/${NAME}
BUILD=.build

all: build

clean:
	rm -rf pkg bin ${BUILD}

build:
	gb build all

docker: docker-build

docker-clean: clean
	docker rmi ${NAME} &>/dev/null || true

docker-build: docker-clean
	docker pull mickep76/alpine-golang:latest
	docker run --rm -it -v "$$PWD":/api -w /api mickep76/alpine-golang:latest
	docker build --pull=true --no-cache -t ${REPO}/${NAME}:${TAG} .
	( cd datavol; docker build --pull=true --no-cache -t ${REPO}/mongo-data:${TAG} . )

docker-push: docker-build
	docker login
	docker push ${REPO}/${NAME}:${TAG}
	docker tag -f ${REPO}/${NAME}:${TAG} ${REPO}/${NAME}:latest
	docker push ${REPO}/${NAME}:latest
	docker tag -f ${REPO}/mongo-data:${TAG} ${REPO}/mongo-data:latest
	docker push ${REPO}/mongo-data:latest
