FROM gliderlabs/alpine

RUN mkdir -p /app/etc

COPY etc/etcdrest.json /app/etc/etcdrest.json
COPY schemas /app/schemas
COPY templates /app/templates

VOLUME ["/app/etcd", "/app/schemas", "/app/templates"]
