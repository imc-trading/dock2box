#!/bin/ash

set -eux

VERS=0.00.1

# Add repositories
cat << EOF >/etc/apk/repositories
http://nl.alpinelinux.org/alpine/v3.4/main
http://nl.alpinelinux.org/alpine/v3.4/community
EOF

# Add pre-reqs
apk update
apk add pciutils \
        parted \
        mdadm \
        lvm2 \
        gptfdisk \
        openssh \
        openssh-client \
        sudo \
        docker \
        curl \
        udev \
        sfdisk \
        iptables \
        ca-certificates \
        tar \
        coreutils \
        jq

# Add initramfs-init script
cp init /usr/share/mkinitfs/initramfs-init

# Add CA certificates
mkdir -p /usr/local/share/ca-certificates/
cp certs/* /usr/local/share/ca-certificates/

# Add SSH keys
mkdir -p /home/dock2box/.ssh
chmod 700 /home/dock2box/.ssh
chown -R dock2box:dock2box /home/dock2box/.ssh
cat ssh/*.rsa >/home/dock2box/.ssh/authorized_keys
chmod 600 /home/dock2box/.ssh/authorized_keys

# Add SSH keys (root temp. workaround)
mkdir -p /root/.ssh
chmod 700 /root/.ssh
chown -R root:root /root/.ssh
cat ssh/*.rsa >/root/.ssh/authorized_keys
chmod 600 /root/.ssh/authorized_keys

# Copy scripts to root
cp functions*.sh /

# Only allow public key login
cat << EOF >>/etc/ssh/sshd_config
PasswordAuthentication no
RSAAuthentication yes
PubkeyAuthentication yes
EOF

echo 'features="network ata base ide scsi usb virtio ext4 dhcp lspci sshd dock2box parted ca-certificates curl docker lvm sgdisk sfdisk udevadm mkfs jq"' >/etc/mkinitfs/mkinitfs.conf
echo "/usr/share/udhcpc/default.script" >>/etc/mkinitfs/features.d/dhcp.files
echo "kernel/net/packet/af_packet.ko" >>/etc/mkinitfs/features.d/dhcp.modules

# bnx2x module is included by default, but depends on crc32c which is not included by default
echo 'kernel/arch/x86/crypto/crc32c-intel.ko' >>/etc/mkinitfs/features.d/network.modules

cat << EOF >/etc/mkinitfs/features.d/lspci.files
/usr/share/hwdata/pci.ids
/usr/sbin/lspci
EOF

# Add motd
rel=$(cat scripts_cache/release)
ts=$(date -u +'%F %T UTC')
sed -e "s/version: x.xx.x/version: ${VERS}/" -e "s/release: xxxxxxx/release: ${rel}/" \
    -e "s/built: xxxx-xx-xx xx:xx:xx xxx/built: ${ts}/" motd >/etc/motd

cat << EOF >/etc/mkinitfs/features.d/sshd.files
/usr/bin/ssh-keygen
/etc/ssh/sshd_config
/usr/sbin/sshd
/usr/lib/ssh/ssh-pkcs11-helper
/etc/motd
EOF

cat << EOF >/etc/mkinitfs/features.d/parted.files
/usr/sbin/partprobe
/usr/sbin/parted
/lib/ld-musl-x86_64.so.1
/usr/lib/libparted.so.2
/usr/lib/libreadline.so.6
/usr/lib/libncursesw.so.6
EOF

cat << EOF >/etc/mkinitfs/features.d/ca-certificates.files
/usr/share/ca-certificates/mozilla
/usr/local/share/ca-certificates
/etc/ca-certificates.conf
/usr/sbin/update-ca-certificates
/etc/ca-certificates/update.d/c_rehash
/usr/bin/c_rehash
EOF

cat << EOF >/etc/mkinitfs/features.d/dock2box.files
/home/dock2box/.ssh/authorized_keys
/root/.ssh/authorized_keys
/functions*.sh
EOF

cat << EOF >/etc/mkinitfs/features.d/docker.files
/usr/bin/docker
/usr/bin/docker-containerd
/usr/bin/docker-runc
/usr/bin/docker-containerd-shim
/usr/bin/docker-containerd-ctr
/sbin/iptables
/usr/lib/libip4tc.so.0
/usr/lib/libip6tc.so.0
/usr/lib/libxtables.so.11
/usr/lib/xtables
/usr/bin/tar
/usr/sbin/chroot
EOF

cat << EOF >/etc/mkinitfs/features.d/docker.modules
kernel/net/bridge
kernel/net/netfilter
kernel/net/ipv4
kernel/net/ipv6/netfilter
kernel/drivers/net/veth.ko
EOF

cat << EOF >/etc/mkinitfs/features.d/curl.files
/usr/bin/curl
/usr/lib/libcurl.so.4
/usr/lib/libssh2.so.1
EOF

cat << EOF >/etc/mkinitfs/features.d/lvm.files
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
/sbin/lvremove
/sbin/pvremove
/lib/libdevmapper-event.so.1.02
EOF

cat << EOF >/etc/mkinitfs/features.d/sgdisk.files
/usr/bin/sgdisk
/lib/ld-musl-x86_64.so.1
/lib/libpopt.so.0
/usr/lib/libstdc++.so.6
/usr/lib/libgcc_s.so.1
EOF

cat << EOF >/etc/mkinitfs/features.d/sfdisk.files
/sbin/sfdisk
/lib/libfdisk.so.1
/lib/libsmartcols.so.1
EOF

cat << EOF >/etc/mkinitfs/features.d/udevadm.files
/sbin/udevadm
EOF

cat << EOF >/etc/mkinitfs/features.d/mkfs.files
/sbin/mkfs.ext4
/lib/libext2fs.so.2
/lib/libcom_err.so.2
/lib/libe2p.so.2
/sbin/mkswap
EOF

cat << EOF >/etc/mkinitfs/features.d/jq.files
/usr/bin/jq
EOF

# Fake uname
mv /bin/uname /bin/uname.orig
cat << EOF >/bin/uname
#!/bin/ash
echo 4.4.17-0-grsec
EOF
chmod +x /bin/uname

# Build initrd and copy kernel
mkinitfs -o initrd
cp /boot/vmlinuz-grsec kernel
