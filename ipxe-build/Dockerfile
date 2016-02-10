FROM gliderlabs/alpine

RUN apk update ;\
    apk add curl gcc make perl musl-dev xz-dev ;\
    mkdir /build

WORKDIR "/build"
ENTRYPOINT [ "make" ]
