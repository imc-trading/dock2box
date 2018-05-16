FROM debian:jessie

RUN set -ex ;\
    apt-get update ;\
    DEBIAN_FRONTEND=noninteractive apt-get install -y -q libguestfs-tools linux-image-3.16.0-4-amd64

COPY build.sh /build.sh
RUN chmod +x /build.sh

WORKDIR /build
ENTRYPOINT /build.sh
