FROM gliderlabs/alpine

RUN apk update ;\
    apk add --no-cache tftp-hpa ;\
    adduser -D tftp

COPY mapfile undionly.kpxe /var/tftpboot/
RUN chmod 0444 /var/tftpboot/*

VOLUME /var/tftpboot
EXPOSE 69/udp

ENTRYPOINT ["in.tftpd", "--foreground"]
CMD ["--verbose", "-B", "1380", "-m", "/var/tftpboot/mapfile", "--user", "tftp", "--secure", "/var/tftpboot/"]
