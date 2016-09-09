#!/bin/ash

set -eux

# Add repositories
cat << EOF >/etc/apk/repositories
http://nl.alpinelinux.org/alpine/v3.4/main
http://nl.alpinelinux.org/alpine/v3.4/community
EOF

# Add pre-reqs
apk update
apk add pciutils

# Configure initrd
(
    cd /usr/share/mkinitfs/
    patch initramfs-init </build/initramfs-init.patch
)

echo 'features="network ata base ide scsi usb virtio ext4 dhcp lspci"' >/etc/mkinitfs/mkinitfs.conf
echo "/usr/share/udhcpc/default.script" >>/etc/mkinitfs/features.d/dhcp.files
echo "kernel/net/packet/af_packet.ko" >>/etc/mkinitfs/features.d/dhcp.modules

#bnx2x module is included by default, but depends on crc32c which is not included by default
echo 'kernel/arch/x86/crypto/crc32c-intel.ko' >>/etc/mkinitfs/features.d/network.modules

cat << EOF >>/etc/mkinitfs/features.d/lspci.files
/usr/share/hwdata/pci.ids
/usr/sbin/lspci
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

# Create package list
cat << EOF >>/etc/apk/world
parted
mdadm
lvm2
gptfdisk
openssh
sudo
docker
curl
udev
bash
sfdisk
EOF

# Add CA certificates
mkdir -p /usr/local/share/ca-certificates/
cp certs/* /usr/local/share/ca-certificates/

# Add startup scripts
cp dock2box /etc/init.d/dock2box
chmod +x /etc/init.d/dock2box

# Add services
rc-update add dock2box
rc-update add sshd
rc-update add hwdrivers

# Add SSH keys
mkdir -p /home/dock2box/.ssh
chmod 700 /home/dock2box/.ssh
chown -R dock2box:dock2box /home/dock2box/.ssh
cat sshkeys/*.rsa >/home/dock2box/.ssh/authorized_keys

# Only allow public key login
cat << EOF >>/etc/ssh/sshd_config
PasswordAuthentication no
RSAAuthentication yes
PubkeyAuthentication yes
EOF

# Lock root account
passwd -d root
passwd -l root

# Add motd
vers=0.00.1
rel=$(cat scripts_cache/release)
ts=$(date -u +'%F %T UTC')
sed -e "s/version: x.xx.x/version: ${vers}/" -e "s/release: xxxxxxx/release: ${rel}/" \
    -e "s/built: xxxx-xx-xx xx:xx:xx xxx/built: ${ts}/" motd >>/etc/motd
sed -e "s/version: x.xx.x/version: ${vers}/" -e "s/release: xxxxxxx/release: ${rel}/" \
    -e "s/built: xxxx-xx-xx xx:xx:xx xxx/built: ${ts}/" -e 's/\\/\\\\/g' motd >>/etc/issue

# Add aliases
cat << EOF >/etc/profile.d/aliases
alias ll="ls -lah"
EOF

# Add scripts
cp scripts_cache/*.sh /root
chmod +x /root/install.sh

# Create apkovl.tar.gz with modules and package list
tar -czf apkovl.tar.gz \
    /lib/modules \
    /etc/apk/world \
    /etc/apk/repositories \
    /etc/sudoers.d/dock2box \
    /etc/passwd /etc/shadow \
    /home/dock2box \
    /usr/local/share/ca-certificates \
    /etc/motd \
    /etc/issue \
    /etc/init.d/dock2box \
    /etc/runlevels/sysinit/dock2box \
    /etc/runlevels/sysinit/sshd \
    /etc/runlevels/sysinit/hwdrivers \
    /etc/profile.d/aliases \
    /root/*.sh
