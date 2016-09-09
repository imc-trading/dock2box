NAME=centos-7.2.1511
ORG=dock2box
VIRT_TAR_OUT=$(ORG)/debian8.5-virt-tar-out
SYS=$(shell uname -s)
REL=$(shell git rev-parse --verify --short HEAD)

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

export:
	docker pull $(VIRT_TAR_OUT):latest
	docker run --privileged --rm -v "$${PWD}":/build --entrypoint /build.sh $(VIRT_TAR_OUT) $(NAME) $(BUILDER)

import: export
	docker build --no-cache --force-rm -t $(ORG)/$(NAME):$(REL) --build-arg NAME=${NAME} .

push: import
	docker login
	docker push ${ORG}/${NAME}:${REL}

latest: push
	docker tag $(ORG)/$(NAME):$(REL) $(ORG)/$(NAME):latest
	docker push $(ORG)/$(NAME):latest
