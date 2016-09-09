KOPTS="vconsole.keymap=us verbose nomodeset crashkernel=auto selinux=0 pcie_aspm=off"

wipe_mdraid() {
    local disk1=$1
    local disk2=$2

    set +e
    # gentoo-ism
    /etc/init.d/mdraid stop

    # Cleanup software RAID
    for MD in $(cd /dev; ls -d md? 2>/dev/null); do
        mdadm --stop ${MD}
        #rm -rf /dev/${MD}
    done

    # Cleanup software RAID
    for MD in $(cd /dev; ls -d md? 2>/dev/null); do
        mdadm --stop ${MD}
    done

    rm -f /usr/lib/udev/rules.d/*-md-*

    # clean up our devices
    udevadm control --reload
    sleep 1
    set -e
}

conf_net() {
    echo ${HOSTNAME} >${ROOT}/etc/hostname
    cat <<EOF >${ROOT}/etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
EOF
    rm -f ${ROOT}/etc/sysconfig/network-scripts/{ifcfg-,routes}*
}

conf_net_dhcp() {
    local device=$1 hwaddr=$2

    cat << EOF >${ROOT}/etc/sysconfig/network-scripts/ifcfg-$device
DEVICE=$device
ONBOOT=yes
HWADDR=$hwaddr
BOOTPROTO=dhcp
DHCP_HOSTNAME=${HOSTNAME%%.*}
EOF
    cat <<EOF >${ROOT}/etc/sysconfig/network
NETWORKING=yes
NETWORKING_IPV6=no
NOZEROCONF=yes
EOF
}

install_grub() {
    local disk=$1 kopts=$2 dracut_kopts

    dracut_kopts=$(chroot ${ROOT} dracut --print-cmdline | xargs printf "%s ")

    cat << EOF >$ROOT/etc/default/grub
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR="Fedora Linux"
GRUB_DEFAULT=saved
GRUB_DISABLE_SUBMENU=true
GRUB_TERMINAL_OUTPUT="console"
GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0 $kopts $dracut_kopts"
GRUB_DISABLE_RECOVERY="true"
EOF

    cat << EOF >$ROOT/boot/grub2/device.map
(hd0)      $disk
EOF

    kver=$(ls $ROOT/lib/modules/ | tail -1) 

    chroot $ROOT grub2-mkconfig -o /boot/grub2/grub.cfg 2>&1
    chroot $ROOT grub2-install $disk
    chroot $ROOT mkinitrd --force /boot/initramfs-${kver}.img $kver 2>&1

#    chroot ${ROOT} grub2-install $disk
#    if [ ${MDRAID-} == 1 ]; then
#        echo "(hd1)      ${SDB}" >>${ROOT}/boot/grub2/device.map
#        chroot ${ROOT} grub2-install ${SDB}
#   fi
}

write_fstab() {
    cat << EOF >${ROOT}/etc/fstab
LABEL=rootfs    /            ext4      defaults         1 1
LABEL=varfs     /var         ext4      defaults         1 1
LABEL=BOOTFS    /boot        vfat      defaults         1 2
LABEL=swapfs    swap         swap      defaults         0 0
EOF
}
