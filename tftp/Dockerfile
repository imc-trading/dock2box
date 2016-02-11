FROM gliderlabs/alpine

ENV IPXE_URL https://git.ipxe.org/ipxe.git/snapshot/e2b1140486e6d5da756d64ae5fc051b79664c6d6.tar.gz
ENV IPXE_VERS e2b1140

ARG SERVER_URI=http://dock2box:8080
ENV SERVER_URI ${SERVER_URI}

RUN echo ${SERVER_URI}

# Install tftp
RUN set -ex ;\
    apk add --no-cache tftp-hpa ;\
    adduser -D tftp

# Compile kpxe and the purge build env.
COPY mapfile /var/tftpboot/
COPY patch-${IPXE_VERS}.txt /patch.txt
COPY embed_template.ipxe /
RUN set -ex ;\
    sed s!__SERVER_URI__!${SERVER_URI}!g /embed_template.ipxe >/embed.ipxe ;\
    cat /embed.ipxe ;\
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
