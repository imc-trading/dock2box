NAME=etcdrest-config
VERSION=0.1
USER=imctrading

all: docker

clean: clean
	docker rmi ${NAME} &>/dev/null || true

docker: clean
	docker build --pull=true --no-cache -t ${USER}/${NAME}:${VERSION} .
	docker tag -f ${USER}/${NAME}:${VERSION} ${USER}/${NAME}:latest

push: docker
	docker login
	docker push ${USER}/${NAME}:${VERSION}
	docker push ${USER}/${NAME}:latest
