NAME=peekaboo
SRCDIR=src/github.com/imc-trading/${NAME}
BUILDDIR=.build
RESDIR=/var/lib/${NAME}
VERSION:=$(shell awk -F '"' '/Version/ {print $$2}' ${SRCDIR}/version.go)
RELEASE:=$(shell date -u +%Y%m%d%H%M)

all: build

clean:
	rm -f *.rpm
	rm -rf pkg bin ${BUILDDIR}

deps:
	go get github.com/constabulary/gb/...

test: clean
	gb test

build: test
	gb build all

update:
	gb vendor update --all

install:
	cp bin/peekaboo /usr/bin
	mkdir -p ${RESDIR}
	cp -r ${SRCDIR}/static ${RESDIR}
	cp -r ${SRCDIR}/templates ${RESDIR}

rpm:
	docker pull mickep76/centos-golang:latest
	docker run --rm -it -v "$$PWD":/go/src/${SRCDIR} -w /go/src/${SRCDIR} mickep76/centos-golang:latest make build-rpm

build-rpm: deps build
	mkdir -p ${BUILDDIR}/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}
	cp -r bin files ${SRCDIR}/templates ${SRCDIR}/static ${BUILDDIR}/SOURCES
	sed -e "s/%NAME%/${NAME}/g" -e "s/%VERSION%/${VERSION}/g" -e "s/%RELEASE%/${RELEASE}/g" \
	${NAME}.spec >${BUILDDIR}/SPECS/${NAME}.spec
	rpmbuild -vv -bb --target="x86_64" --clean --define "_topdir $$(pwd)/${BUILDDIR}" ${BUILDDIR}/SPECS/${NAME}.spec
	mv ${BUILDDIR}/RPMS/x86_64/*.rpm .
