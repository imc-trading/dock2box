USER=imctrading
NAME=dock2box-ipxe-build

all: build

clean:
	docker rmi ${NAME} &>/dev/null || true

build: clean
	docker build --pull=true --no-cache -t ${USER}/${NAME}:latest .

push: build
	docker login
	docker push ${USER}/${NAME}:latest
