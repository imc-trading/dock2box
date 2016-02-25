FROM gliderlabs/alpine

RUN mkdir -p /app/etc

COPY etc/etcdrest.yaml /app/etc/etcdrest.yaml
COPY schemas /app/schemas
COPY templates /app/templates
COPY images /app/images
COPY scripts /app/scripts

RUN set -ex ;\
    apk add --no-cache --virtual build perl-digest-sha1 ;\
    for f in $(find /app/scripts -type f); do \
        sha256sum $f >$f.sha256 ;\
    done ;\
    apk del build

VOLUME ["/app/etcd", "/app/schemas", "/app/templates", "/app/images", "/app/scripts"]
