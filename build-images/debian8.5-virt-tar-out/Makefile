NAME=debian8.5-virt-tar-out
ORG=dock2box

all: push

clean:
	docker rmi ${NAME} &>/dev/null || true

build: clean										   
	docker build --pull=true --no-cache -t ${ORG}/${NAME}:latest .

push: build
	docker login
	docker push ${ORG}/${NAME}:latest
