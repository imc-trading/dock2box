SYS=$(shell uname -s)
REL=$(shell git rev-parse --verify --short HEAD)
VIRT_TAR_OUT=dock2box/debian8.5-virt-tar-out
SETUPIMG=$(VIRT_TAR_OUT):latest

#BUILDER=virtualbox-iso

ifeq ($(SYS), Darwin)
BUILDER=vmware-iso
else
BUILDER=qemu
endif

# display video output when run by humans
export HEADLESS = false

all: latest

# but not when run headless
headless: export HEADLESS=true
headless: latest

clean:
	rm -rf packer_output Dockerfile

test:
	@type packer;
	@type docker;
	docker ps >/dev/null
	packer validate packer.json

build: clean test
	packer build -only=$(BUILDER) -var "vm_name=$(NAME)" packer.json

import: build
	docker pull $(SETUPIMG) &>/dev/null
	docker run --privileged --rm -v "$${TEAMCITY_BUILD_CHECKOUTDIR:-$$PWD}":/build --entrypoint /build.sh $(SETUPIMG) $(NAME) $(BUILDER) | docker import - $(REG)/$(REPO)/$(NAME):$(REL)

docker-build: import
	@echo 'FROM $(REG)/$(REPO)/$(NAME):$(REL)' > ./Dockerfile
	docker build .
