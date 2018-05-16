FROM alpine

RUN mkdir -p /var/lib/dock2box/static/boot

COPY kernel /var/lib/dock2box/static/boot/kernel
COPY initrd /var/lib/dock2box/static/boot/initrd
COPY img /var/lib/dock2box/static/boot/img

RUN find /var/lib/dock2box/static -type f -exec chmod 644 {} \; ;\
    find /var/lib/dock2box/static -type d -exec chmod 755 {} \;

VOLUME ["/var/lib/dock2box/static/boot"]
ENTRYPOINT ["tail", "-f", "/dev/null"]
