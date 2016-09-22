NAME=alpine3.4.3
REPO=dock2box
REG=
SYS=$(shell uname -s)
REL=$(shell git rev-parse --verify --short HEAD)
VIRT_TAR_OUT=dock2box/debian8.5-virt-tar-out

# To use VM-Ware builder, specify:
# make BUILDER=vmware-iso

ifeq ($(SYS), Darwin)
BUILDER=virtualbox-iso
else
BUILDER=qemu
endif

all: latest

clean:
	rm -rf packer_output

check:
	@type packer;
	@type docker;
	docker ps >/dev/null

test: check
	packer validate template.json

build: clean test
	packer build -only=$(BUILDER) -var "vm_name=$(NAME)" template.json

import: build
	docker pull $(REG)$(VIRT_TAR_OUT):latest
	docker run --privileged --rm -v "$${PWD}":/build --entrypoint /build.sh $(REG)$(VIRT_TAR_OUT):latest $(NAME) $(BUILDER) | docker import - $(REG)$(REPO)/$(NAME):$(REL)

push: import
	docker login
	docker push $(REG)$(REPO)/$(NAME):$(REL)

latest: push
	docker tag $(REG)$(REPO)/$(NAME):$(REL) $(REG)$(REPO)/$(NAME):latest
	docker push $(REG)$(REPO)/$(NAME):latest
