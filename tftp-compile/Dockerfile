FROM gliderlabs/alpine

ENV IPXE_URL https://git.ipxe.org/ipxe.git/snapshot/e2b1140486e6d5da756d64ae5fc051b79664c6d6.tar.gz
ENV IPXE_VERS e2b1140

# Install tftp
RUN set -ex ;\
    apk add --no-cache tftp-hpa bash ;\
    adduser -D tftp

# Compile kpxe and the purge build env.
COPY mapfile /var/tftpboot/
COPY patch.txt embed.ipxe embed_template.ipxe start.sh /
RUN set -ex ;\
    apk add --no-cache --virtual build curl gcc make perl musl-dev xz-dev ;\
    curl -o ipxe.tar.gz ${IPXE_URL} ;\
    tar zxf ipxe.tar.gz ;\
    cd /ipxe-${IPXE_VERS} && patch -p1 < /patch.txt ;\
    cd /ipxe-${IPXE_VERS}/src && make CFLAGS_config="-DIMAGE_COMBOOT=on" EMBED=/embed.ipxe bin/undionly.kpxe ;\
    cp /ipxe-${IPXE_VERS}/src/bin/undionly.kpxe /var/tftpboot/undionly.kpxe ;\
    chmod 0444 /var/tftpboot/* ;\
    rm -f /patch.txt /ipxe.tar.gz ;\
    chmod 755 /start.sh

VOLUME /var/tftpboot
EXPOSE 69/udp

ENTRYPOINT ["/start.sh"]
CMD ["${SERVER_NAME}"]
