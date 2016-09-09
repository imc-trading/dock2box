FROM scratch

ARG NAME=?
ENV NAME ${NAME}

ADD packer_output/${NAME}.tar.gz /

COPY dracut.conf /etc/
RUN rm -f /etc/dracut.conf.d/ami.conf

CMD /bin/bash
