NAME=boot-image
BOOTIMG=alpine3.4.3
USER=mickep76
REPO=dock2box
REG=
REL=$(shell git rev-parse --verify --short HEAD)
TS=$(shell date -u +'%F %T UTC')

all: push

clean:
	rm -f initrd kernel
	docker rmi ${NAME} &>/dev/null || true

build:
	docker pull $(REG)$(REPO)/$(BOOTIMG):latest 
	docker run --rm -v "$${PWD}":/build -w /build $(REG)$(REPO)/$(BOOTIMG):latest ./setup.sh $(REL) "$(TS)"
	docker build --pull=true --no-cache -t $(REG)$(REPO)/$(NAME):$(REL) --build-arg REL=$(REL) --build-arg TS="$(TS)" .
	docker tag $(REG)$(REPO)/$(NAME):$(REL) $(REEG)$(REPO)/$(NAME):latest

push: build
	docker login -u $(USER)
	docker push $(REG)$(REPO)/$(NAME):$(REL)
	docker push $(REG)$(REPO)/$(NAME):latest
