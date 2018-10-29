ARG src="/go/src/github.com/imc-trading/dock2box"
ARG rel=?
ARG codename="Dapper Snapper"
#ARG onload="201805"

#
# Agent Builder
#
FROM mickep76/alpine-golang AS agent-builder

#COPY certs /usr/local/share/ca-certificates
#RUN update-ca-certificates

ARG src
WORKDIR $src

COPY agent agent/
COPY models models/
COPY common common/
COPY vendor vendor/
COPY glide.yaml glide.yaml
COPY glide.lock glide.lock

RUN glide up
RUN cd agent && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

#
# Initrd Builder
#
FROM alpine:3.7 AS initrd-builder

COPY boot/repositories /etc/apk/repositories

RUN set -eux ;\
    apk update ;\
    apk add \
        kbd-bkeymaps \
        busybox-initscripts \
        e2fsprogs-libs \
        e2fsprogs \
        pv \
        ipmitool \
        alpine-sdk \
        linux-vanilla \
        linux-vanilla-dev \
        libtool \
        libpcap-dev \
        autoconf \
        automake \
        musl-dev \
        pciutils \
        parted \
        mdadm \
        lvm2 \
        gptfdisk \
        openssh \
        openssh-client \
        sudo \
        curl \
        udev \
        sfdisk \
        iptables \
        ca-certificates \
        tar \
        coreutils \
        jq \
        cmake \
        v86d \
        ncurses \
        util-linux \
        xfsprogs \
        btrfs-progs \
        mdadm \
        util-linux \
        bc \
        git \
        bash \
        sgdisk \
        open-vm-tools \
        grub \
        openssl \
        device-mapper

RUN set -eux ;\
    grep -v 'root:' /usr/share/mkinitfs/passwd >/tmp/passwd ;\
    echo 'root:$6$BaHYi8W6U.pAapmr$AfxQe39FlxKZh1EiFbNeROwiUWadGFUbgHs.ZZK0RNi8M4giXhDDnX/SukpZRXelKPt8B3ZJHcVAsHD.qroTw1:0:0:root:/root:/bin/ash' >/usr/share/mkinitfs/passwd ;\
    cat /tmp/passwd >>/usr/share/mkinitfs/passwd

COPY boot/init.sh /usr/share/mkinitfs/initramfs-init
RUN chmod +x /usr/share/mkinitfs/initramfs-init

COPY certs /usr/local/share/ca-certificates
RUN update-ca-certificates

ARG src
COPY --from=agent-builder $src/agent/agent /usr/bin/agent

RUN mkdir -p /var/log

#RUN git clone "https://github.com/imc-trading/dock2box.git" /etc/dock2box/scripts

RUN set -eux ;\
    mkdir -p /root/.ssh ;\
    chmod 700 /root/.ssh ;\
    chown -R root:root /root/.ssh ;\
    touch /root/.ssh/authorized_keys

COPY boot/sshd_config /etc/ssh/sshd_config

COPY boot/sshkeys /tmp/sshkeys
RUN [ -n "$(ls /tmp/sshkeys)" ] && cat tmp/sshkeys/*.rsa >>/root/.ssh/authorized_keys ;\
    chmod 600 /root/.ssh/authorized_keys ;\
    rm -rf /tmp/sshkeys

ARG rel
ARG codename
RUN set -eux ;\
    echo "Dock2Box Linux release ${rel} (${codename})" >/etc/redhat-release ;\
    cp /etc/redhat-release /etc/motd

COPY boot/dock2box-800x600.ppm /etc/dock2box.ppm

RUN echo "options uvesafb mode_option=800x600-32 scroll=ywrap" >/etc/modprobe.d/uvesafb.conf

ARG onload
RUN curl -O https://www.openonload.org/download/openonload-${onload}.tgz ;\
    tar zxf openonload-${onload}.tgz ;\
    cd openonload-${onload}/scripts ;\
    ./onload_build --kernelver $(ls /lib/modules | grep -v hardened | tail -1) --kernel ;\
    ./onload_install --kernelver $(ls /lib/modules | grep -v hardened | tail -1) --nobuild --kernelfiles --modprobe --initscript

COPY boot/features.d /etc/mkinitfs/features.d

#RUN set -eux ;\
#    echo -n 'features="' >/etc/mkinitfs/mkinitfs.conf ;\
#    ls /etc/mkinitfs/features.d | sed -e 's/\.files//g' -e 's/\.modules//g' | uniq | tr '\n' ' '>>/etc/mkinitfs/mkinitfs.conf ;\
#    echo '"' >>/etc/mkinitfs/mkinitfs.conf

RUN echo 'features="keymap dhcp_ext network_ext network ata base ide scsi usb virtio ext4 dhcp lspci sshd dock2box parted ca-certificates tar chroot curl lvm gdisk sgdisk sfdisk udevadm mkfs jq v86d pv ncurses tput setterm xfsprogs xfs vmxnet mdadm agent wipefs bc git onload ipmitool nvme grub dmsetup ahci"' >/etc/mkinitfs/mkinitfs.conf

RUN mkinitfs -o /initrd $(ls /lib/modules | grep -v hardened | tail -1)

#
# Boot
#
FROM alpine

RUN mkdir -p /var/lib/dock2box/static/boot

COPY --from=initrd-builder /boot/vmlinuz /var/lib/dock2box/static/boot/kernel
COPY --from=initrd-builder /initrd /var/lib/dock2box/static/boot/initrd
COPY boot/img /var/lib/dock2box/static/boot/img

RUN find /var/lib/dock2box/static -type f -exec chmod 644 {} \; ;\
    find /var/lib/dock2box/static -type d -exec chmod 755 {} \;

VOLUME ["/var/lib/dock2box/static/boot"]
ENTRYPOINT ["tail", "-f", "/dev/null"]
