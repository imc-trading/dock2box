FROM gliderlabs/alpine

COPY bin/d2bsrv /d2bsrv
COPY schemas /schemas
COPY templates /templates
COPY img /img

EXPOSE 8080
ENTRYPOINT ["/d2bsrv"]
