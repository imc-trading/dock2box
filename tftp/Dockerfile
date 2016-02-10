FROM gliderlabs/alpine

ENV IPXE_URL https://git.ipxe.org/ipxe.git/snapshot/e2b1140486e6d5da756d64ae5fc051b79664c6d6.tar.gz
ENV IPXE_VERS e2b1140

# Install tftp
RUN set -ex ;\
    apk add --no-cache tftp-hpa ;\
    adduser -D tftp

# Compile kpxe and the purge build env.
COPY mapfile ipxelinux.ipxe /var/tftpboot/
RUN set -ex ;\
    apk add --no-cache --virtual build curl gcc make perl musl-dev xz-dev ;\
    curl -o ipxe.tar.gz ${IPXE_URL} ;\
    tar zxf ipxe.tar.gz ;\
    ( cd ipxe-${IPXE_VERS}/src && make CFLAGS_config="-DIMAGE_COMBOOT=on -DCONSOLE_CMD=on" EMBED=/var/tftpboot/ipxelinux.ipxe bin/undionly.kpxe ) ;\
    cp /ipxe-${IPXE_VERS}/src/bin/undionly.kpxe /var/tftpboot/undionly.kpxe ;\
    rm -rf /ipxe-${IPXE_VERS} ;\
    apk del build

RUN chmod 0444 /var/tftpboot/*

VOLUME /var/tftpboot
EXPOSE 69/udp

ENTRYPOINT ["in.tftpd", "--foreground"]
CMD ["--verbose", "-B", "1380", "-m", "/var/tftpboot/mapfile", "--user", "tftp", "--secure", "/var/tftpboot/"]
