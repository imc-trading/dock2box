FROM alpine:3.7

ENV IPXE_URL https://git.ipxe.org/ipxe.git/snapshot/fbe8c52d0d9cdb3d6f5fe8be8edab54618becc1f.tar.gz
ENV IPXE_VERS fbe8c52

# Install tftp
RUN set -ex ;\
    apk add --no-cache tftp-hpa ;\
    adduser -D tftp

# Compile kpxe and the purge build env.
COPY mapfile /var/tftpboot/
COPY patch-${IPXE_VERS}.txt /patch.txt
COPY embed.ipxe /
RUN set -ex ;\
    apk add --no-cache --virtual build curl gcc make perl musl-dev xz-dev ;\
    curl -o ipxe.tar.gz ${IPXE_URL} ;\
    tar zxf ipxe.tar.gz ;\
    cd /ipxe-${IPXE_VERS} && patch -p1 < /patch.txt ;\
    cd /ipxe-${IPXE_VERS}/src && make CFLAGS_config="-DIMAGE_COMBOOT=on" EMBED=/embed.ipxe bin/undionly.kpxe ;\
    cp /ipxe-${IPXE_VERS}/src/bin/undionly.kpxe /var/tftpboot/undionly.kpxe ;\
    chmod 0444 /var/tftpboot/* ;\
    rm -rf /ipxe-${IPXE_VERS} ;\
    apk del build ;\
    rm -f /patch.txt /embed_template.ipxe /embed.ipxe /ipxe.tar.gz

VOLUME /var/tftpboot
EXPOSE 69/udp

ENTRYPOINT ["in.tftpd", "--foreground"]
CMD ["--verbose", "-B", "1380", "-m", "/var/tftpboot/mapfile", "--user", "tftp", "--secure", "/var/tftpboot/"]
