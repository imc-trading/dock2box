#!/bin/sh

on_exit_shell() {
	echo "Failed! Entering emergency shell"
	echo "Type exit to reboot"
	sh
	reboot -d 3 -f
}

trap on_exit_shell EXIT

TITLE="\e[0;37m"
HELP="\e[1;30m"
INFO="\e[1;30m"
WARN="\e[0;33m"
FATAL="\e[0;31m"

CLEAR="\033[0m"
CLEAR_LINE="\033[K"

COLS=$(tput cols)
LINES=$(tput lines)

/bin/busybox mkdir -p /usr/bin /usr/sbin /proc /sys /dev \
	/media/cdrom /media/usb /tmp /run
/bin/busybox --install -s

# Basic environment
export PATH=/usr/bin:/bin:/usr/sbin:/sbin

# Needed devs
[ -c /dev/null ] || mknod -m 666 /dev/null c 1 3

# Basic mounts
mount -t proc -o noexec,nosuid,nodev proc /proc
mount -t sysfs -o noexec,nosuid,nodev sysfs /sys

info() {
	if [ -e /tmp/progress_bar ]; then
		clear_lines $(($LINES - 6)) >/dev/tty2
		print_center $(($LINES - 6)) "$1" $INFO >/dev/tty2
		[ -n "${2:-}" ] && echo $2 >/tmp/progress_bar
	fi
	printf "${INFO}$1${CLEAR}\n"
}

warn() {
	if [ -e /tmp/progress_bar ]; then
		clear_lines $(($LINES - 6)) >/dev/tty2
		print_center $(($LINES - 6)) "$1" $WARN >/dev/tty2
		[ -n "${2:-}" ] && echo $2 >/tmp/progress_bar
	fi
	printf "${WARN}WARN: $1${CLEAR}\n"
}

fatal() {
	if [ -e /tmp/progress_bar ]; then
		clear_lines $(($LINES - 6)) >/dev/tty2
		print_center $(($LINES - 6)) "$1" $FATAL >/dev/tty2
		echo 100 >/tmp/progress_bar
	fi
	printf "${FATAL}FATAL: $1${CLEAR}\n"
	exit 1
}

print_center() {
	local line=$1 text="$2" color="${3:-$INFO}" center

	center=$(($COLS / 2 - ${#text} / 2))
	tput cup $line $center
	echo -en "${color}${text}${CLEAR}"
}

print_pos() {
	local line=$1 col=$2 text="$3" color="${4:-$INFO}"

	tput cup $line $col
	echo -en "${color}${text}${CLEAR}"
}

clear_lines() {
	local beg=$1 num=${2:-0}

	for line in $(seq 0 $num); do
		tput cup $(($beg + $line)) 0
		echo -en ${CLEAR_LINE}
	done
}

get_interface() {
	local hwaddr="$1"

	ip addr show | grep -i -B 1 $hwaddr | head -1 | cut -d : -f 2 | tr -d '[[:space:]]'
}

ip_choose_if() {
	local hwaddr=$1

	NET_OK="false"
	if [ -n "$hwaddr" ]; then
		PROTO="dhcp"
		INTERFACE=$(get_interface $hwaddr)
		HWADDR=$hwaddr

		ifconfig $INTERFACE 0.0.0.0 >&2
		udhcpc -i $INTERFACE -f -q -F ${HWADDR//:/} -n >&2

		[ $? -eq 0 ] && NET_OK="true" && return
	fi

	for i in /sys/class/net/e*; do
		if [ -e "$i" ]; then
			PROTO="dhcp"
			INTERFACE=${i##*/}
			HWADDR=$(cat /sys/class/net/$INTERFACE/address)

			ifconfig $INTERFACE 0.0.0.0 >&2
			udhcpc -i $INTERFACE -f -q -F ${HWADDR//:/} -n >&2

			[ $? -eq 0 ] && NET_OK="true" && return
		fi
	done
}

ip_up_if() {
	for i in /sys/class/net/e*; do
		if [ -e "$i" ]; then
			INTERFACE=${i##*/}
			ifconfig $INTERFACE up >&2
			#			CARRIER=$(cat /sys/class/net/$INTERFACE/carrier)
			#			[ ${CARRIER:-0} -eq 1 ] && ifconfig $INTERFACE up >&2
		fi
	done
}

get_first_disk() {
	[ -b "/dev/vda" ] && echo "/dev/vda" && return
	[ -b "/dev/sda" ] && echo "/dev/sda"
}

wait_disk() {
	local disk=$1

	while [ ! -b $disk ]; do
		sleep 1
		echo -n .
	done
}

# Get kernel options
set -- $(cat /proc/cmdline)

myopts="mac quiet splash video debug dma modules usbdelay blacklist ntp_servers server elastic_url"

for opt; do
	for i in $myopts; do
		case "$opt" in
		$i=*) eval "KOPT_${i}='${opt#*=}'" ;;
		$i) eval "KOPT_${i}=true" ;;
		false$i) eval "KOPT_${i}=false" ;;
		esac
	done
done

# Enable debug
[ -n "$KOPT_debug" ] && set -x

info "Setting US keymap"
zcat /usr/share/bkeymaps/us/us.bmap.gz | loadkmap

# No DMA
[ "$KOPT_dma" == "false" ] && modprobe libata dma=0

# Hide Kernel messages
[ "$KOPT_quiet" == "true" ] && dmesg -n 1

# Blacklist modules
for i in ${KOPT_blacklist/,/ }; do
	echo "blacklist $i" >>/etc/modprobe.d/boot-opt-blacklist.conf
done

# Setup /dev
mount -t devtmpfs -o exec,nosuid,mode=0755,size=2M devtmpfs /dev 2>/dev/null ||
	mount -t tmpfs -o exec,nosuid,mode=0755,size=2M tmpfs /dev
[ -d /dev/pts ] || mkdir -m 755 /dev/pts
[ -c /dev/ptmx ] || mknod -m 666 /dev/ptmx c 5 2
# Make sure /dev/null is setup correctly
[ -f /dev/null ] && rm -f /dev/null
[ -c /dev/null ] || mknod -m 666 /dev/null c 1 3
mount -t devpts -o gid=5,mode=0620,noexec,nosuid devpts /dev/pts
[ -d /dev/shm ] || mkdir /dev/shm
mount -t tmpfs -o nodev,nosuid,noexec shm /dev/shm

# Start 4 Virtual Consoles
cat <<EOF >/etc/profile
TTY=\$(tty)
PS1="\h [\${TTY##/dev/}]# "
EOF

openvt -c 2
openvt -c 3
openvt -c 4
openvt -c 5
openvt -c 6

# Load VESA and framebuffer
if [ -n "${KOPT_video:-}" ]; then
	info "Loading VESA driver and framebuffer"
	echo "options uvesafb mode_option=${KOPT_video} scroll=ywrap" >/etc/modprobe.d/uvesafb.conf
	modprobe uvesafb
	modprobe fbcon

	COLS=$(tput cols)
	LINES=$(tput lines)
fi

if [ "$KOPT_splash" == "true" ]; then
	chvt 2
	setterm -msg off -cursor off -foreground green >/dev/tty2
	tput clear >/dev/tty2
	cat <<EOF >/etc/fbsplash.conf
IMAGE_ALIGN=CM

BAR_LEFT=50
BAR_TOP=525

BAR_WIDTH=700
BAR_HEIGHT=20

BAR_R=35
BAR_G=35
BAR_B=35
EOF

	mkfifo /tmp/progress_bar
	fbsplash -s /etc/dock2box.ppm -i /etc/fbsplash.conf -f /tmp/progress_bar &
	print_center 1 "Dock2Box" $TITLE >/dev/tty2
	print_pos 2 1 "Switch Virtual Terminal using ALT+F[1-6]" $HELP >/dev/tty2
fi

# Load drivers
info "Loading boot drivers" 10

modprobe -a $(echo "$KOPT_modules" | tr ',' ' ') 2>/dev/null
if [ -f /etc/modules ]; then
	sed 's/\#.*//g' </etc/modules |
		while read module args; do
			modprobe -q $module $args
		done
fi

# Loading network drivers
info "Loading network drivers" 20
nlplug-findfs -p /sbin/mdev ${KOPT_debug:+-d} \
	${KOPT_usbdelay:+-t $(($KOPT_usbdelay * 1000))}

#info "Loading onload drivers" 25
#modprobe -a onload

# Load IPMI
info "Loading IPMI drivers" 28
modprobe ipmi_devintf
modprobe ipmi_si

# Get network device
info "Configure network" 30
ip_choose_if $KOPT_mac

[ "$NET_OK" == "true" ] || fatal "Failed to configure network"
HWADDR=$(cat /sys/class/net/$INTERFACE/address)

# Enable all interfaces with a carrier for LLDP
ip_up_if

# Get hostname
info "Set hostname ${HWADDR//:/}" 45
hostname ${HWADDR//:/}

ifconfig lo up

# Update title with ip
ip=$(ip address show $INTERFACE | awk '/inet / {print $2}' | awk -F / '{print $1}')
print_center 4 "HWAddr: ${HWADDR} IP: ${ip}" $HELP >/dev/tty2

# Update title with sn
sn=$(cat /sys/devices/virtual/dmi/id/product_serial)
print_center 5 "SN: ${sn}" $HELP >/dev/tty2

info "Loading hardware drivers" 65
find /sys -name modalias -type f -print0 | xargs -0 sort -u |
	xargs modprobe -b -a 2>/dev/null

# Run twice so we detect all devices
find /sys -name modalias -type f -print0 | xargs -0 sort -u |
	xargs modprobe -b -a 2>/dev/null

# Start sshd
info "Starting sshd" 70
mkdir -p /var/empty
ssh-keygen -A
/usr/sbin/sshd

# Update CA certificates
info "Update CA certificates" 80
mkdir -p /etc/ssl/certs /etc/ca-certificates/update.d
update-ca-certificates

[ "$KOPT_quiet" == "true" ] || cat /etc/motd

# Wait for disk to become available
disk=$(get_first_disk)
wait_disk $disk

# Create symlink to mtab
ln -s /proc/mounts /etc/mtab

# Set time
info "Sync time" 90
[ -z "$KOPT_ntp_servers" ] && KOPT_ntp_servers="ntp1 ntp2 ntp3"
ntpd -d -q -n -p $KOPT_ntp_servers

# Update scripts
info "Update scripts" 95
(
	cd /etc/dock2box/scripts
	git pull
)
mkdir -p /var/log

info "Start agent and wait..." 100
[ -z "$KOPT_server" ] && KOPT_server="dock2box"
/usr/bin/agent --endpoints $KOPT_server
on_exit_shell
