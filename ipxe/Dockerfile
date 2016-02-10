FROM gliderlabs/alpine

RUN apk update ;\
    apk add --no-cache tftp-hpa ;\
    adduser -D tftp

COPY undionly.kpxe /var/tftpboot/undionly.kpxe
RUN chmod 0444 /var/tftpboot/undionly.kpxe

VOLUME /var/tftpboot
EXPOSE 69/udp

ENTRYPOINT ["in.tftpd", "--foreground"]
CMD ["--verbose", "--user", "tftp", "--secure", "/var/tftpboot/"]
