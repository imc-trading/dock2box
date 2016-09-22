KOPTS="vconsole.keymap=us verbose nomodeset crashkernel=auto selinux=0 pcie_aspm=off"

conf_net() {
    echo $HOSTNAME >$SYSROOT/etc/hostname
    cat <<EOF >$SYSROOT/etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
EOF
    rm -f $SYSROOT/etc/sysconfig/network-scripts/{ifcfg-,routes}*
}

conf_net_dhcp() {
    cat << EOF >$SYSROOT/etc/sysconfig/network-scripts/ifcfg-$INTERFACE
DEVICE=$INTERFACE
ONBOOT=yes
HWADDR=$HWADDR
BOOTPROTO=$PROTO
DHCP_HOSTNAME=${HOSTNAME%%.*}
EOF
    cat <<EOF >$SYSROOT/etc/sysconfig/network
NETWORKING=yes
NETWORKING_IPV6=no
NOZEROCONF=yes
EOF
}

conf_net_static() {
    cat << EOF >$SYSROOT/etc/sysconfig/network-scripts/ifcfg-$INTERFACE
DEVICE=$INTERFACE
ONBOOT=yes
HWADDR=$HWADDR
BOOTPROTO=$PROTO
IPADDR=$IP
NETMASK=$NETMASK
EOF
    cat << EOF >$SYSROOT/etc/sysconfig/network
NETWORKING=yes
NETWORKING_IPV6=no
GATEWAY=$GW
NOZEROCONF=yes
EOF

    cat << EOF >$SYSROOT/etc/resolv.conf
search $SEARCH
nameserver $DNS1
nameserver $DNS2
EOF
}

install_grub() {
    local disk=$1 kopts="$2"

    cat << EOF >$SYSROOT/etc/default/grub
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR="Fedora Linux"
GRUB_DEFAULT=saved
GRUB_DISABLE_SUBMENU=true
GRUB_TERMINAL_OUTPUT="console"
GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0 $kopts"
GRUB_DISABLE_RECOVERY="true"
EOF

    cat << EOF >$SYSROOT/boot/grub2/device.map
(hd0)      $disk
EOF

    chroot $SYSROOT grub2-mkconfig -o /boot/grub2/grub.cfg 2>&1
    chroot $SYSROOT grub2-install $disk
}

write_fstab() {
    local sysroot=$1

    cat << EOF >$SYSROOT/etc/fstab
LABEL=rootfs    /            ext4      defaults         1 1
LABEL=varfs     /var         ext4      defaults         1 1
LABEL=BOOTFS    /boot        vfat      defaults         1 2
LABEL=swapfs    swap         swap      defaults         0 0
EOF
}
