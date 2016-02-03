FROM gliderlabs/alpine

RUN mkdir -p /app/etc

COPY etc/etcdrest.yaml /app/etc/etcdrest.yaml
COPY schemas /app/schemas
COPY templates /app/templates

VOLUME ["/app/etcd", "/app/schemas", "/app/templates"]
