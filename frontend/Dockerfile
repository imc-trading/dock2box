FROM alpine

RUN mkdir -p /var/lib/dock2box
COPY app /var/lib/dock2box/static

VOLUME ["/var/lib/dock2box/static"]
ENTRYPOINT ["tail", "-f", "/dev/null"]
