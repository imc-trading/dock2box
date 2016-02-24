FROM gliderlabs/alpine

RUN mkdir -p /app/etc

COPY etc/etcdrest.yaml /app/etc/etcdrest.yaml
COPY schemas /app/schemas
COPY templates /app/templates
COPY images /app/images
COPY scripts /app/scripts

VOLUME ["/app/etcd", "/app/schemas", "/app/templates", "/app/images", "/app/scripts"]
