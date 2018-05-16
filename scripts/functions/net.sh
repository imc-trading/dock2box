to_dec() {
	local hex=$(echo $1 | tr 'a-z' 'A-Z')

	echo "ibase=16; $hex" | bc
}

net_if_primary() {
	route | awk '/0.0.0.0\s+UG/ {print $8}'
}

net_if_lldp_primary() {
	for f in /sys/class/net/e*; do
		i=$(basename $f)
		lldpctl $i | awk -F ':' -v i=$i '/PortDescr:/ && $4=="pri" {print i}'
	done
}

net_if_lldp_vlan() {
	lldpctl $i | awk -F '[:,]\ +' '/VLAN:/ {print $2}'
}

net_if_hwaddr() {
	local interface=$1

	cat /sys/class/net/${interface}/address
}

net_if_name() {
	local hwaddr=$1

	grep -il $hwaddr /sys/class/net/*/address | awk -F '/' '{print $5}'
}

net_if_bus() {
	local interface=$1

	readlink -f /sys/class/net/$interface | awk -F '/' '{print $6}' | awk -F ':' '{print $2}'
}

net_if_slot() {
	local interface=$1

	readlink -f /sys/class/net/$interface | awk -F '/' '{print $6}' | awk -F ':' '{print $3}' | awk -F '.' '{print $1}'
}

net_if_func() {
	local interface=$1

	readlink -f /sys/class/net/$interface | awk -F '/' '{print $6}' | awk -F ':' '{print $3}' | awk -F '.' '{print $2}'
}

net_if_dev_port() {
	local interface=$1

	if [ -e "/sys/class/net/$interface/dev_port" ]; then
		cat /sys/class/net/$interface/dev_port
	fi
}

net_if_name_path() {
	local hwaddr=$1

	local name=$(net_if_name $hwaddr)
	local bus=$(net_if_bus $name)
	local slot=$(net_if_slot $name)
	local func=$(net_if_func $name)
	local dev_port=$(net_if_dev_port $name)

	local bus_dec=$(to_dec $bus)
	local slot_dec=$(to_dec $slot)
	local func_dec=$(to_dec $func)
	local dev_port_dec=$(to_dec $dev_port)

	if [ "$dev_port" != "0" ]; then
		echo enp${bus_dec}s${slot_dec}f${func_dec}d${dev_port_dec}
		return
	fi

	if [ "$func" != "0" ]; then
		echo enp${bus_dec}s${slot_dec}f${func_dec}
		return
	fi

	if [ -e "/sys/bus/pci/devices/0000:${bus}:${slot}.1" ]; then

		if [ -e "/sys/bus/pci/devices/0000:${bus}:${slot}.1/dev_port" ]; then
			local sec_dev_port=$(cat /sys/bus/pci/devices/0000:${bus}:${slot}.1/dev_port)
			if [ "$sec_dev_port" != "0" ]; then
				echo enp${bus_dec}s${slot_dec}f${func_dec}d${dev_port_dec}
				return
			fi
		fi

		echo enp${bus_dec}s${slot_dec}f${func_dec}
		return
	fi

	echo enp${bus_dec}s${slot_dec}
}

net_cleanup_network_scripts() {
	local root=$1

	rm -f $(ls $root/etc/sysconfig/network-scripts/ifcfg-* | grep -v lo)
	rm -f $root/etc/sysconfig/network-scripts/routes*
}

net_set_etc_hostname() {
	local root=$1 hostname=$2

	echo $hostname >$root/etc/hostname
}

net_create_etc_hosts() {
	local root=$1

	cat <<EOF >$root/etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
EOF
}

net_create_sysconfig_network() {
	local root=$1

	cat <<EOF >$root/etc/sysconfig/network
NETWORKING=yes
NETWORKING_IPV6=no
NOZEROCONF=yes
EOF
}

net_set_udev_name_path() {
	local root=$1

	[ -d $root/etc/udev/rules.d ] && mkdir -p $root/etc/udev/rules.d

	cat <<EOF >$root/etc/udev/rules.d/80-net-name-path.rules
ACTION!="add", GOTO="net_name_slot_end"
SUBSYSTEM!="net", GOTO="net_name_slot_end"
NAME!="", GOTO="net_name_slot_end"

IMPORT{cmdline}="net.ifnames"
ENV{net.ifnames}=="0", GOTO="net_name_slot_end"

#NAME=="", ENV{ID_NET_NAME_ONBOARD}!="", NAME="\$env{ID_NET_NAME_ONBOARD}"
#NAME=="", ENV{ID_NET_NAME_SLOT}!="", NAME="\$env{ID_NET_NAME_SLOT}"
NAME=="", ENV{ID_NET_NAME_PATH}!="", NAME="\$env{ID_NET_NAME_PATH}"

LABEL="net_name_slot_end"
EOF
}

net_dhcp() {
	local root=$1 hostname=$2 hwaddr=$3

	[ -e $root/etc/sysconfig/network-scripts/ifcfg-eth0 ] && rm -f $root/etc/sysconfig/network-scripts/ifcfg-eth0

	name=$(net_if_name_path $hwaddr)

	cat <<EOF >$root/etc/sysconfig/network-scripts/ifcfg-$name
DEVICE=$name
ONBOOT=yes
HWADDR=$hwaddr
BOOTPROTO=dhcp
DHCP_HOSTNAME=${hostname%%.*}

#BOOTPROTO=static
#IPADDR=
#NETMASK=
EOF
}

net_etc_resolv_conf() {
	local root=$1 search=$2 dns_servers="$3" dns_options="${4:-}"

	cat <<EOF >$root/etc/resolv.conf
options $dns_options
domain $search
search $search
EOF

	for srv in $dns_servers; do
		echo "nameserver $srv" >>${root}/etc/resolv.conf
	done
}

net_static() {
	local root=$1 ip=$2 netmask=$3 gw=$4 hwaddr=$5

	[ -e $root/etc/sysconfig/network-scripts/ifcfg-eth0 ] && rm -f $root/etc/sysconfig/network-scripts/ifcfg-eth0

	name=$(net_if_name_path $hwaddr)

	cat <<EOF >$root/etc/sysconfig/network-scripts/ifcfg-$name
DEVICE=$name
ONBOOT=yes
BOOTPROTO=static
IPADDR=$ip
NETMASK=$netmask
EOF

	cat <<EOF >>${root}/etc/sysconfig/network-scripts/ifcfg-$name
HWADDR=$hwaddr
EOF

	cat <<EOF >>$root/etc/sysconfig/network
GATEWAY=$gw
EOF
}
