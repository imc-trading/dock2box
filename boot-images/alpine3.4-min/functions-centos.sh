KOPTS="vconsole.keymap=us verbose nomodeset crashkernel=auto selinux=0 pcie_aspm=off"

conf_net() {
    local hostname=$1 sysroot=$2

    echo $hostname >$sysroot/etc/hostname
    cat <<EOF >${ROOT}/etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
EOF
    rm -f $sysroot/etc/sysconfig/network-scripts/{ifcfg-,routes}*
}

conf_net_dhcp() {
    local device=$1 hwaddr=$2 hostname=$3 sysroot=$4

    cat << EOF >$sysroot/etc/sysconfig/network-scripts/ifcfg-$device
DEVICE=$device
ONBOOT=yes
HWADDR=$hwaddr
BOOTPROTO=dhcp
DHCP_HOSTNAME=${hostname%%.*}
EOF
    cat <<EOF >${ROOT}/etc/sysconfig/network
NETWORKING=yes
NETWORKING_IPV6=no
NOZEROCONF=yes
EOF
}

install_grub() {
    local disk=$1 sysroot=$2 kopts="$3"

    cat << EOF >$sysroot/etc/default/grub
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR="CentOS Linux"
GRUB_DEFAULT=saved
GRUB_DISABLE_SUBMENU=true
GRUB_TERMINAL_OUTPUT="console"
GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0 $kopts"
GRUB_DISABLE_RECOVERY="true"
EOF

    cat << EOF >$sysroot/boot/grub2/device.map
(hd0)      $disk
EOF

    chroot $sysroot grub2-mkconfig -o /boot/grub2/grub.cfg 2>&1
    chroot $sysroot grub2-install $disk
}

write_fstab() {
    local sysroot=$1

    cat << EOF >$sysroot/etc/fstab
LABEL=rootfs    /            ext4      defaults         1 1
LABEL=varfs     /var         ext4      defaults         1 1
LABEL=BOOTFS    /boot        vfat      defaults         1 2
LABEL=swapfs    swap         swap      defaults         0 0
EOF
}