#!/bin/ash

set -eu
set -x

REL=$1
CODENAME="$2"
BUILD_ROOT=$(pwd)

#
# Add applications and configuration
#

# Add repositories
cat <<EOF >/etc/apk/repositories
http://dl-cdn.alpinelinux.org/alpine/v3.7/main
http://dl-cdn.alpinelinux.org/alpine/v3.7/community
EOF

# Add pre-reqs
apk update
apk add pv \
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

# Configure Go
export GOROOT="/usr/lib/go"
export GOPATH="/go"
export PATH="${GOPATH}/bin:$PATH"
mkdir -p "${GOPATH}/src" "${GOPATH}/bin"

# Setup root passwd
grep -v 'root:' /usr/share/mkinitfs/passwd >/tmp/passwd
echo 'root:$6$BaHYi8W6U.pAapmr$AfxQe39FlxKZh1EiFbNeROwiUWadGFUbgHs.ZZK0RNi8M4giXhDDnX/SukpZRXelKPt8B3ZJHcVAsHD.qroTw1:0:0:root:/root:/bin/ash' >/usr/share/mkinitfs/passwd
cat /tmp/passwd >>/usr/share/mkinitfs/passwd

# Add initramfs-init script
cp init.sh /usr/share/mkinitfs/initramfs-init

# Add CA certificates
#mkdir -p /usr/local/share/ca-certificates/
#[ -n "$(ls .certs-cache)" ] && cp .certs-cache/* /usr/local/share/ca-certificates/
#update-ca-certificates

# Add client
cp ./client /usr/bin/client
mkdir -p /etc/dock2box

# Add SSH keys
mkdir -p /root/.ssh
chmod 700 /root/.ssh
chown -R root:root /root/.ssh
touch /root/.ssh/authorized_keys
[ -n "$(ls sshkeys)" ] && cat sshkeys/*.rsa >>/root/.ssh/authorized_keys
chmod 600 /root/.ssh/authorized_keys

# Setup SSH config
mkdir -p /etc/ssh
touch /etc/ssh/sshd_config
cat <<EOF >/etc/ssh/sshd_config
PasswordAuthentication yes
PermitRootLogin yes
RSAAuthentication yes
PubkeyAuthentication yes
EOF

KVER=$(ls /lib/modules | grep -v hardened | tail -1)
echo "Kernel Version: $KVER" >&2

# Onload driver
#ONLOAD_VER="201606-u1.3"
#(
#	curl -O http://www.openonload.org/download/openonload-${ONLOAD_VER}.tgz
#	tar zxf openonload-${ONLOAD_VER}.tgz
#	cd openonload-${ONLOAD_VER}/scripts
#	./onload_build --kernelver $KVER --kernel
#	./onload_install --kernelver $KVER --nobuild --kernelfiles --modprobe --initscript
#)

# Add motd and redhat-release
echo "Dock2Box Linux release ${REL} (${CODENAME})" >/etc/redhat-release
cp /etc/redhat-release /etc/motd

# Copy splash screen
cp dock2box-800x600.ppm /etc/dock2box.ppm

# Default uvesafb mode
echo "options uvesafb mode_option=800x600-32 scroll=ywrap" >/etc/modprobe.d/uvesafb.conf

#
# Files to include in initrd
#

echo 'features="network ata base ide scsi usb virtio ext4 dhcp lspci sshd dock2box parted ca-certificates tar chroot curl lvm gdisk sgdisk sfdisk udevadm mkfs jq v86d pv ncurses tput setterm xfsprogs xfs vmxnet mdadm agent wipefs bc git onload ipmitool nvme grub dmsetup"' >/etc/mkinitfs/mkinitfs.conf
echo "/usr/share/udhcpc/default.script" >>/etc/mkinitfs/features.d/dhcp.files
echo "kernel/net/packet/af_packet.ko" >>/etc/mkinitfs/features.d/dhcp.modules

# bnx2x module is included by default, but depends on crc32c which is not included by default
echo 'kernel/arch/x86/crypto/crc32c-intel.ko' >>/etc/mkinitfs/features.d/network.modules

cat <<EOF >>/etc/mkinitfs/features.d/base.files
/etc/keymap/us.bmap.gz
EOF

cat <<EOF >/etc/mkinitfs/features.d/lspci.files
/usr/share/hwdata/pci.ids
/usr/sbin/lspci
EOF

cat <<EOF >/etc/mkinitfs/features.d/sshd.files
/usr/bin/ssh-keygen
/etc/ssh/sshd_config
/usr/sbin/sshd
/usr/lib/ssh/ssh-pkcs11-helper
/etc/motd
/usr/bin/scp
EOF

cat <<EOF >/etc/mkinitfs/features.d/parted.files
/usr/sbin/partprobe
/usr/sbin/parted
/lib/ld-musl-x86_64.so.1
/usr/lib/libparted.so.2
/usr/lib/libreadline.so.6
/usr/lib/libncursesw.so.6
EOF

cat <<EOF >/etc/mkinitfs/features.d/ca-certificates.files
/usr/share/ca-certificates/mozilla
/usr/local/share/ca-certificates
/etc/ca-certificates.conf
/usr/sbin/update-ca-certificates
/etc/ca-certificates/update.d/c_rehash
/usr/bin/c_rehash
EOF

cat <<EOF >/etc/mkinitfs/features.d/dock2box.files
/root/.ssh/authorized_keys
/etc/dock2box.ppm
/etc/redhat-release
/opt/sbin/lldpcli
/opt/sbin/lldpctl
/opt/sbin/lldpd
/usr/sbin/lldpcli
/usr/sbin/lldpctl
/usr/sbin/lldpd
/opt/lib/liblldpctl.a
/opt/lib/liblldpctl.la
/opt/lib/liblldpctl.so
/opt/lib/liblldpctl.so.4
/opt/lib/liblldpctl.so.4.8.0
/opt/lib/pkgconfig/lldpctl.pc
/opt/var/run/.saveme
EOF

cat <<EOF >/etc/mkinitfs/features.d/chroot.files
/usr/sbin/chroot
EOF

cat <<EOF >/etc/mkinitfs/features.d/tar.files
/usr/bin/tar
EOF

cat <<EOF >/etc/mkinitfs/features.d/curl.files
/usr/bin/curl
/usr/lib/libcurl.so.4
/usr/lib/libssh2.so.1
EOF

cat <<EOF >/etc/mkinitfs/features.d/lvm.files
/sbin/vgs
/sbin/lvs
/sbin/pvs
/sbin/vgscan
/sbin/lvscan
/sbin/pvscan
/sbin/vgcreate
/sbin/lvcreate
/sbin/pvcreate
/sbin/vgremove
/sbin/vgreduce
/sbin/lvremove
/sbin/pvremove
/sbin/pvchange
/sbin/lvchange
/sbin/vgchange
/lib/libdevmapper-event.so.1.02
EOF

cat <<EOF >/etc/mkinitfs/features.d/sgdisk.files
/usr/bin/sgdisk
/lib/ld-musl-x86_64.so.1
/lib/libpopt.so.0
/usr/lib/libstdc++.so.6
/usr/lib/libgcc_s.so.1
EOF

cat <<EOF >/etc/mkinitfs/features.d/gdisk.files
/usr/bin/gdisk
/lib/libuuid.so.1
EOF

cat <<EOF >/etc/mkinitfs/features.d/sfdisk.files
/sbin/sfdisk
/lib/libfdisk.so.1
/lib/libsmartcols.so.1
EOF

cat <<EOF >/etc/mkinitfs/features.d/udevadm.files
/sbin/udevadm
EOF

cat <<EOF >/etc/mkinitfs/features.d/mkfs.files
/sbin/mkfs.ext4
/lib/libext2fs.so.2
/lib/libcom_err.so.2
/lib/libe2p.so.2
/sbin/mkswap
EOF

cat <<EOF >/etc/mkinitfs/features.d/jq.files
/usr/bin/jq
EOF

cat <<EOF >/etc/mkinitfs/features.d/ttylog.files
/usr/sbin/ttylog
EOF

cat <<EOF >/etc/mkinitfs/features.d/v86d.files
/sbin/v86d
/etc/modprobe.d/uvesafb.conf
EOF

cat <<EOF >/etc/mkinitfs/features.d/v86d.modules
kernel/drivers/video/console/fbcon.ko
kernel/drivers/video/fbdev/uvesafb.ko
EOF

cat <<EOF >/etc/mkinitfs/features.d/pv.files
/usr/bin/pv
EOF

cat <<EOF >/etc/mkinitfs/features.d/ncurses.files
/usr/lib/libncursesw.so.6
/etc/terminfo
EOF

cat <<EOF >/etc/mkinitfs/features.d/tput.files
/usr/bin/tput
EOF

cat <<EOF >/etc/mkinitfs/features.d/setterm.files
/usr/bin/setterm
EOF

cat <<EOF >/etc/mkinitfs/features.d/xfsprogs.files
/sbin/xfs_repair
/sbin/fsck.xfs
/sbin/mkfs.xfs
/lib/libuuid.so.1
/lib/ld-musl-x86_64.so.1
/lib/libblkid.so.1
EOF

cat <<EOF >/etc/mkinitfs/features.d/xfs.modules
kernel/fs/xfs
kernel/fs/xfs/xfs.ko
EOF

cat <<EOF >/etc/mkinitfs/features.d/vmxnet.modules
kernel/drivers/net/vmxnet3
kernel/drivers/net/vmxnet3/vmxnet3.ko
misc/vmxnet.ko
EOF

cat <<EOF >/etc/mkinitfs/features.d/mdadm.modules
kernel/drivers/md/raid0.ko
kernel/drivers/md/raid10.ko
kernel/drivers/md/raid1.ko
kernel/drivers/md/raid456.ko
EOF

cat <<EOF >/etc/mkinitfs/features.d/mdadm.files
etc/mdadm.conf
etc/conf.d/mdadm
lib/udev/rules.d/63-md-raid-arrays.rules
lib/udev/rules.d/64-md-raid-assembly.rules
sbin/mdadm
EOF

cat <<EOF >/etc/mkinitfs/features.d/agent.files
/usr/bin/agent
/etc/dock2box
/var/log/dock2box
/var/lib/dock2box
EOF

cat <<EOF >/etc/mkinitfs/features.d/wipefs.files
/sbin/wipefs
/lib/ld-musl-x86_64.so.1
/lib/libblkid.so.1
/lib/ld-musl-x86_64.so.1
/lib/libuuid.so.1
EOF

cat <<EOF >/etc/mkinitfs/features.d/bc.files
/usr/bin/bc
EOF

cat <<EOF >/etc/mkinitfs/features.d/git.files
/usr/bin/git
/usr/libexec/git-core/git
/usr/libexec/git-core/git-clone
/usr/libexec/git-core/git-pull
/usr/libexec/git-core/git-remote-https
EOF

cat <<EOF >/etc/mkinitfs/features.d/onload.modules
extra
EOF

cat <<EOF >/etc/mkinitfs/features.d/onload.files
/etc/modprobe.d/onload.conf
/etc/depmod.d/onload.conf
EOF

cat <<EOF >/etc/mkinitfs/features.d/ipmitool.files
/usr/sbin/ipmitool
/usr/share/ipmitool/oem_ibm_sel_map
/usr/lib/libreadline.so.7.0
/usr/lib/libreadline.so.7
EOF

cat <<EOF >/etc/mkinitfs/features.d/ipmitool.modules
kernel/drivers/char/ipmi
EOF

cat <<EOF >/etc/mkinitfs/features.d/nvme.modules
kernel/drivers/nvme
EOF

cat <<EOF >/etc/mkinitfs/features.d/grub.files
/usr/sbin/grub-probe
/lib/libdevmapper.so.1.02
/usr/sbin/grub-mkconfig
/usr/sbin/grub-install
/usr/bin/grub-editenv
/usr/bin/grub-file
/usr/bin/grub-script-check
/usr/sbin/grub-set-default
/usr/sbin/grub-mkconfig
/usr/share/grub/grub-mkconfig_lib
EOF

cat <<EOF >/etc/mkinitfs/features.d/dmsetup.files
/sbin/dmsetup
EOF

# Build initrd and copy kernel
mkinitfs -o initrd $KVER
cp /boot/vmlinuz kernel
