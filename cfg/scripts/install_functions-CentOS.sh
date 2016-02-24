# Centos-7 Version of wipe_mdraid()
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
    #rm -rf /dev/md

    # Cleanup software RAID
    for MD in $(cd /dev; ls -d md? 2>/dev/null); do
        mdadm --stop ${MD}
    done

    # Redhat-ism
    rm -f /usr/lib/udev/rules.d/*-md-*

    # clean up our devices
    udevadm control --reload
    sleep 1
    set -e
}

# Centos-7 Version: Inject /etc/hosts
write_net_config_common() {
    echo ${HOSTNAME} >${ROOT}/etc/hostname
    cat <<EOF >${ROOT}/etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
EOF
    rm -f ${ROOT}/etc/sysconfig/network-scripts/{ifcfg-,routes}*
}

# Centos-7 Version: Inject net config (static)
write_net_config_static() {
    DNS="..."
    local INTERFACES=$(get_config_subnodes "host/${HOSTNAME}/interface")
    local ROUTES=$(get_config_subnodes "host/${HOSTNAME}/route")

    echo "${IP} ${HOSTNAME} ${HOSTNAME_SHORT}" >>${ROOT}/etc/hosts

    if [ -n "${INTERFACES}" ]; then
        if [ -f ${ROOT}/etc/udev/rules.d/70-persistent-net.rules ]; then
            rm ${ROOT}/etc/udev/rules.d/70-persistent-net.rules
        fi
        # This systems's interfaces seem registered in etcd, let's use this info
        DONE=""
        for IF_DEVICE in ${INTERFACES}; do
            IF_HWADDR=$(get_config_val "host/${HOSTNAME}/interface/${IF_DEVICE}/hwaddr")
            IF_IP=$(get_config_val "host/${HOSTNAME}/interface/${IF_DEVICE}/ip")
            IF_NETMASK=$(get_config_val "host/${HOSTNAME}/interface/${IF_DEVICE}/netmask")
            # if default, handle as special case for device name (map to real device, e.g. "eth0")
            if [[ "${IF_DEVICE}" == "default" ]]; then
                # we only have the default interface defined in etcd, no eth0, etc
                IF_DEVICE=$(grep --with-filename ${IF_HWADDR} /sys/class/net/*/address | cut -d/ -f5 | cut -d. -f1 | sort -u)
            fi
            # If we haven't already configured this interface (could happen if, e.g., both default and eth0 are defined)
            if [[ ${DONE} != *"${IF_DEVICE}"* ]]; then
                IF_CFG=${ROOT}/etc/sysconfig/network-scripts/ifcfg-${IF_DEVICE}
                if [[ ${IF_DEVICE} == [a-z0-9]*.[0-9]* ]]; then
                    echo "VLAN=yes" > ${IF_CFG}
                fi
                if [ -n "${IF_IP}" ]; then
                    echo "DEVICE=${IF_DEVICE}"   >> ${IF_CFG}
                    echo "ONBOOT=yes"            >> ${IF_CFG}
                    echo "BOOTPROTO=static"      >> ${IF_CFG}
                    echo "HWADDR=${IF_HWADDR}"   >> ${IF_CFG}
                    echo "IPADDR=${IF_IP}"       >> ${IF_CFG}
                    echo "NETMASK=${IF_NETMASK}" >> ${IF_CFG}
                else
                    echo "DEVICE=${IF_DEVICE}"   >> ${IF_CFG}
                    echo "ONBOOT=no"             >> ${IF_CFG}
                    echo "HWADDR=${IF_HWADDR}"   >> ${IF_CFG}
                fi
                # Add udev-rule to maintain consistent names going forward:
                echo SUBSYSTEM==\"net\", ACTION==\"add\", ATTR{address}==\"${IF_HWADDR}\", KERNEL==\"eth*\", NAME=\"${IF_DEVICE}\" >>${ROOT}/etc/udev/rules.d/70-persistent-net.rules
                DONE="${DONE} ${IF_DEVICE}"
            fi
        done
    fi

    for ROUTE in ${ROUTES}; do
        ROUTE_GW=$(get_config_val "host/${HOSTNAME}/route/${ROUTE}/gw")
        ROUTE_DEV=$(get_config_val "host/${HOSTNAME}/route/${ROUTE}/interface")
        echo "${ROUTE} via ${ROUTE_GW}" >> ${ROOT}/etc/sysconfig/network-scripts/route-${ROUTE_DEV}
    done

    DOMAIN=$(get_config_val "global/domain")
    SITE=$(get_config_val "host/${HOSTNAME}/site")

    cat <<EOF >${ROOT}/etc/sysconfig/network
NETWORKING=yes
NETWORKING_IPV6=no
GATEWAY=${GATEWAY}
NOZEROCONF=yes
EOF
    # generate /etc/resolv.conf using nearest DNS servers
    echo "options timeout:2 attempts:1 ndots:1" > ${ROOT}/etc/resolv.conf
    echo "domain ${DOMAIN}" >> ${ROOT}/etc/resolv.conf
    echo "search ${SITE}.${DOMAIN} ${DOMAIN}" >> ${ROOT}/etc/resolv.conf
    # Detect the 2 closest DNS servers
    counter=0
    while ! ping -c 1 -W 1 ${GATEWAY} &>/dev/null; do
        # This is an incredibly obscure, non-googlable construct.  
        # Please rewrite this or document what it does and how, as 3 of us could not make head nor tail of it:
	    "write_net_config_static, ${counter}. attempt."
        sleep 1
        echo -n .
        counter=$((counter+1))
        if [[ "${counter}" -gt 100 ]]; then
            echo "The network has failed to come up after 100 seconds.  Giving up"
            exit 1
        fi
    done; echo "up"
    # don't trust this to be passed as a global env variable:
    DNS="..."
    fping -C1 -a ${DNS} 2>&1 | awk '/^[0-9].*[0-9]$/{ print $1, $3 }' | sort -nk2 | awk '{ if (NR < 3 ) print "nameserver", $1}' >> ${ROOT}/etc/resolv.conf
}

# Centos-7 Version: Inject net config (dhcp)
write_net_config_dhcp() {
    cat <<EOF >${ROOT}/etc/sysconfig/network-scripts/ifcfg-${BOOT_ETH}
DEVICE=${BOOT_ETH}
ONBOOT=yes
HWADDR=${BOOT_MAC}
BOOTPROTO=dhcp
DHCP_HOSTNAME=${HOSTNAME_SHORT}
EOF
    cat <<EOF >${ROOT}/etc/sysconfig/network
NETWORKING=yes
NETWORKING_IPV6=no
NOZEROCONF=yes
EOF
}

# Centos-7 Version: Figure out the complete list of boot args
get_boot_args() {
    local customargs="${1-}"
    local autoargs=$(chroot ${ROOT} dracut --print-cmdline | xargs printf "%s ")
    
    # FIXME this should ensure that overrides/duplicates are merged
    echo "${autoargs} ${KOPTS_DEFAULT} ${customargs}"
}

# Centos-7 Version: Inject bootloader config
grub_config_and_install() {
    local bootargs=${1}

    # Default GRUB config
    cat <<EOF >${ROOT}/etc/default/grub
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR="D2B Linux"
GRUB_DEFAULT=saved
GRUB_DISABLE_SUBMENU=true
GRUB_TERMINAL_OUTPUT="console"
GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0 ${bootargs}"
GRUB_DISABLE_RECOVERY="true"
EOF

    # GRUB device map
    cat <<EOF >${ROOT}/boot/grub2/device.map
# this device map was generated by anaconda
(hd0)      ${SDA}
EOF
    
    # Generate GRUB config
    chroot ${ROOT} grub2-mkconfig -o /boot/grub2/grub.cfg 2>&1 | grep -v lvmetad

    # Install bootloader to disk
    chroot ${ROOT} grub2-install ${SDA}
    if [ ${MDRAID-} == 1 ]; then
        echo "(hd1)      ${SDB}" >>${ROOT}/boot/grub2/device.map
        chroot ${ROOT} grub2-install ${SDB}
   fi
}

# Centos-7 Version: Inject fstab
write_fstab() {
    cat <<EOF >${ROOT}/etc/fstab
LABEL=rootfs    /            ext4      defaults         1 1
LABEL=fs        /fs          ext4      defaults         1 1
LABEL=log       /fs/log      ext4      defaults         1 1
LABEL=BOOTFS    /boot        vfat      defaults         1 2
LABEL=swapfs    swap         swap      defaults         0 0
EOF
}

