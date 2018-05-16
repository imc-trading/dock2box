#!/bin/sh

# Change to vt1 console
chvt 1

set -eux

export PATH=/usr/bin:/bin:/usr/sbin:/sbin

OUT=$(mktemp)
ROOT=/mnt/sysimage

fatal() {
	echo "FATAL: $1" >&2
	exit 1
}

d2b_stage() {
	echo "STAGE:$2: $1" >&2
}

# Verify input
if [ -n "${IP:-}" ]; then
	[ -n "${NETMASK:-}" ] || fatal "NETMASK env. variable isn't set"
	[ -n "${GW:-}" ] || fatal "GW env. variable isn't set"
	[ -n "${DNS_SERVERS:-}" ] || fatal "DNS_SERVERS env. variable isn't set"
fi

[ -n "${IMAGE:-}" ] || fatal "IMAGE env. variable isn't set"

# Set defaults
[ -n "${DISTRO:-}" ] || DISTRO=centos
[ -n "${GPT:-}" ] || GPT=true
[ -n "${REGISTRY:-}" ] || REGISTRY=registry:8080
[ -n "${TAG:-}" ] || TAG=latest
[ -n "${ROOT_SIZE:-}" ] || ROOT_SIZE=10G
[ -n "${SWAP_SIZE:-}" ] || SWAP_SIZE=4G
[ -n "${DOWNLOAD_SIZE:-}" ] || DOWNLOAD_SIZE=10G
[ -n "${VOLUMES:-}" ] || VOLUMES="varfs:vg_root:lv_var:/var:ext4:10G homefs:vg_root:lv_home:/home:ext4:5G"
[ -n "${KOPTS:-}" ] || KOPTS=""
[ -n "${DNS_SEARCH:-}" ] || DNS_SEARCH="example.com"
[ -n "${DNS_OPTIONS:-}" ] || DNS_OPTIONS=""
[ -n "${MDRAID:-}" ] || MDRAID="false"

# Load functions
for f in $(find /etc/dock2box/scripts/functions -type f); do
	source $f
done

# Load distribution specific functions
[ -f /etc/dock2box/scripts/functions-${DISTRO}.sh ] || fatal "Distribution is not supported"
source /etc/dock2box/scripts/functions-${DISTRO}.sh

#
# Disk
#

# Unmount filesystems
d2b_stage "Prepare provisioning" 5
fs_umount_recurse $ROOT

# Get disks
disk=$(disk_get_first)
disk2=$(disk_get_second)

# Wipe MBR
d2b_stage "Wipe disks" 10
disk_mbr_wipe $disk
if [ -n "${disk2:-}" ]; then
	disk_mbr_wipe $disk2
else
	MDRAID="false"
fi

# Wipe MDRaid
mdraid_wipe

# Wipe LVM
lvm_wipe $disk
[ "${MDRAID:-}" == "true" ] && lvm_wipe $disk2

# Wipe disks
disk_part_wipe $disk
[ "${MDRAID:-}" == "true" ] && disk_part_wipe $disk2

# Partition disks
d2b_stage "Setup disks" 15
if [ "${GPT:-}" == "true" ]; then
	disk_part_gpt $disk
	[ "${MDRAID:-}" == "true" ] && disk_part_gpt $disk2
else
	disk_part_mbr $disk
	[ "${MDRAID:-}" == "true" ] && disk_part_mbr $disk2
fi

disk_part=$(disk_part_name $disk)
grubdev=${disk}
bootdev=${disk_part}2
rootdev=${disk_part}3

if [ "${MDRAID:-}" == "true" ]; then
	mdraid_create $disk $disk2

	KOPTS="${KOPTS} rd.md.uuid=$(mdadm --detail /dev/md0 | awk '/UUID/ { print $3 }')"
	KOPTS="${KOPTS} rd.md.uuid=$(mdadm --detail /dev/md1 | awk '/UUID/ { print $3 }')"
	KOPTS="${KOPTS} rd.dm=0"

	grubdev="/dev/sda /dev/sdb"
	bootdev="/dev/md0"
	rootdev="/dev/md1"
fi

# Create boot partition
disk_part_boot $bootdev

# Create lvm physical volume
lvm_create_pv $rootdev

# Create lvm volume group
lvm_create_vg $rootdev vg_root

# Create logical volumes
lvm_create_lv vg_root lv_root ${ROOT_SIZE:-10G}
lvm_create_lv vg_root lv_swap ${SWAP_SIZE:-4G}
lvm_create_lv vg_root lv_download ${DOWNLOAD_SIZE:-10G}

# Create filesystems
fs_create ext4 /dev/vg_root/lv_root rootfs
fs_create swap /dev/vg_root/lv_swap swapfs
fs_create ext4 /dev/vg_root/lv_download downloadfs

# Mount root filesystem
fs_mount ext4 /dev/vg_root/lv_root $ROOT
fs_mount vfat $bootdev $ROOT/boot
fs_mount ext4 /dev/vg_root/lv_download $ROOT/download

# Mount dev, proc and sys filesystems
fs_mount bind /dev $ROOT/dev
fs_mount bind /dev/pts $ROOT/dev/pts
fs_mount bind /proc $ROOT/proc
fs_mount bind /sys $ROOT/sys

# Create fstab
fstab_create $ROOT

# Create MD-Raid config
[ "${MDRAID:-}" == "true" ] && mdraid_config $ROOT

# Create lv, fs and mount for volumes
for vol in $VOLUMES; do
	# $1 (name), $2 (vg), $3 (lv), $4 (mount), $5 (fs), $6 (size)
	IFS=: && set -- $vol && unset IFS
	lvm_create_lv $2 $3 $6
	fs_create $5 /dev/$2/$3 $1
	fs_mount $5 /dev/$2/$3 $ROOT/$4
	fstab_add $ROOT LABEL=$1 $4 $5 defaults 1 1
done

#
# Docker
#

# Pull docker image
d2b_stage "Download image" 30
docker_pull $REGISTRY $IMAGE $TAG $ROOT/download

#
# Network
#

# Install bootloader
d2b_stage "Install bootloader" 55
install_grub $ROOT "${grubdev}" "${KOPTS:-}" "${MDRAID}"

# Cleanup network scripts
d2b_stage "Setup networking" 75
net_cleanup_network_scripts $ROOT

# Get primary interface
if [ -z "${PRIMARY_IF:-}" ]; then
	interface=$(net_if_primary)
	hwaddr=$(net_if_hwaddr $interface)
else
	hwaddr=$PRIMARY_IF
	interface=$(net_if_name $hwaddr)
fi

# Create ifrename config
echo ${RENAME_IF:-} >$ROOT/etc/ifrename.conf

# Setup udev rules
net_set_udev_name_path $ROOT

# Set hostname
net_set_etc_hostname $ROOT $HOSTNAME

# Create /etc/hosts
net_create_etc_hosts $ROOT

# Create /etc/sysconfig/network
net_create_sysconfig_network $ROOT

# Setup network
if [ -z "${IP:-}" ]; then
	net_dhcp $ROOT $HOSTNAME $hwaddr
else
	net_static $ROOT $IP $NETMASK $GW $hwaddr
	net_etc_resolv_conf $ROOT $DNS_SEARCH "$DNS_SERVERS" "$DNS_OPTIONS"
fi

# Cleanup ssh keys
d2b_stage "Cleanup" 80
rm -f $ROOT/etc/ssh/ssh_host_*key*

# Unmount dev, proc and sys filesystems
fs_umount_recurse $ROOT

# Remove tmp volume
lvm_remove_lv vg_root lv_download

# Sync
disk_sync

# Finished
d2b_stage "Finished provisioning, waiting for reboot" 100
