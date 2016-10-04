FROM nginx

MAINTAINER Michael Persson, <michael.ake.persson@gmail.com>

ARG REL=?
ARG TS=?

ENV REL ${REL}
ENV TS ${TS}

# Set motd
RUN printf "dock2box, release: ${REL} (${TS})\n" >>/etc/motd
RUN echo '[ ! -z "$TERM" -a -r /etc/motd ] && cat /etc/motd' >>/etc/bash.bashrc

# Copy config and content
COPY nginx.conf /etc/nginx/nginx.conf
COPY kernel /usr/share/nginx/html/kernel
COPY initrd /usr/share/nginx/html/initrd

# Fix permissions
RUN set -ex ;\
    /bin/bash -c "find /usr/share/nginx/html -type f -exec chmod 644 {} \;" ;\
    /bin/bash -c "find /usr/share/nginx/html -type d -exec chmod 755 {} \;"

# Generate self-signed SSL certificate
ENV SSL_DIR /etc/nginx/ssl
ENV SSL_CN dock2box
ENV SSL_O dock2box
ENV SSL_C NL

RUN set -ex ;\
    mkdir -p ${SSL_DIR} ;\
    openssl req -nodes -new -x509 -keyout ${SSL_DIR}/server.key -out ${SSL_DIR}/server.crt -subj "/CN=${SSL_CN}/O=${SSL_O}/C=${SSL_C}"

EXPOSE 80 443
