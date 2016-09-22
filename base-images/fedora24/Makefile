NAME=fedora24
REPO=dock2box
REG=
SYS=$(shell uname -s)
REL=$(shell git rev-parse --verify --short HEAD)
VIRT_TAR_OUT=debian8.5-virt-tar-out

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

test:
	@type packer;
	@type docker;
	docker ps >/dev/null
	packer validate template.json

build: clean test
	packer build -only=$(BUILDER) -var "vm_name=$(NAME)" template.json

import: build
	docker pull $(REG)$(VIRT_TAR_OUT):latest
	docker run --privileged --rm -v "$${PWD}":/build --entrypoint /build.sh $(REG)$(VIRT_TAR_OUT):latest $(NAME) $(BUILDER) | docker import -c "LABEL dock2box.type='host-image'" - $(REG)$(REPO)/$(NAME):$(REL)

push: import
	docker login
	docker push $(REG)$(REPO)/$(NAME):$(REL)

latest: push
	docker tag $(REG)$(REPO)/$(NAME):$(REL) $(REG)$(REPO)/$(NAME):latest
	docker push $(REG)$(REPO)/$(NAME):latest
