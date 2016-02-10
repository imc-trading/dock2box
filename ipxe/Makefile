USER=imctrading
NAME=dock2box-ipxe
IPXE_URL=https://git.ipxe.org/ipxe.git/snapshot/e2b1140486e6d5da756d64ae5fc051b79664c6d6.tar.gz
IPXE_VERS=e2b1140

all: build

clean:
	docker rmi ${NAME} &>/dev/null || true
	rm -rf ipxe.tar.gz ipxe-${IPXE_VERS} undionly.kpxe
                                                                                               
build: #clean                                                                                   
#	docker run --rm=true -it -v "$$PWD":/build -w /build ${USER}/${NAME}-build:latest ipxe-build
	docker build --pull=true --no-cache -t ${USER}/${NAME}:latest .

push: build
	docker login
	docker push ${USER}/${NAME}:latest

ipxe-build:
	curl -o ipxe.tar.gz ${IPXE_URL}
	tar zxf ipxe.tar.gz 
	cd ipxe-${IPXE_VERS}/src && make CFLAGS_config="-DIMAGE_COMBOOT=on -DCONSOLE_CMD=on" EMBED=../../ipxelinux.ipxe bin/undionly.kpxe
	cp ipxe-${IPXE_VERS}/src/bin/undionly.kpxe undionly.kpxe
