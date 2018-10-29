FROM alpine:3.7

RUN apk update && apk add cdrkit syslinux curl p7zip

RUN mkdir -p /iso/syslinux && cp /usr/share/syslinux/isolinux.bin \
    /usr/share/syslinux/ldlinux.c32 \
    /usr/share/syslinux/vesamenu.c32 \
    /usr/share/syslinux/libcom32.c32 \
    /usr/share/syslinux/libutil.c32 \
    /usr/share/syslinux/chain.c32 \
    /iso/syslinux

COPY isolinux.cfg /iso/syslinux
COPY dock2box.png /iso/syslinux

ARG rel=?
RUN echo "MENU TITLE Dock2Box Boot ISO ($rel)" >>/iso/syslinux/isolinux.cfg

RUN mkdir -p /iso/d2b
COPY --from=<image>/d2b/boot /var/lib/dock2box/static/boot/kernel /iso/d2b/vmlinuz
COPY --from=<image>/d2b/boot /var/lib/dock2box/static/boot/initrd /iso/d2b/initrd.gz

RUN mkisofs -o bootiso.iso \
   -b syslinux/isolinux.bin \
   -c syslinux/boot.cat \
   -no-emul-boot \
   -boot-load-size 4 \
   -boot-info-table \
   iso
