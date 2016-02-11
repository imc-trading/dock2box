USER=imctrading
NAME=dock2box-tftp
SERVER_URI=http://dock2box:8080

all: build

clean:
	docker rmi ${NAME} &>/dev/null || true

build: clean                                                                                   
	docker build --pull=true --no-cache -t ${USER}/${NAME}:latest --build-arg SERVER_URI=${SERVER_URI} .

push: build
	docker login
	docker push ${USER}/${NAME}:latest
