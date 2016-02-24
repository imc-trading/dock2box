#!/bin/bash

# Uses the following options passed by PXE menu
#
# D2B_DEBUG             flag
# D2B_DISTRO            string
# D2B_BOOT_HWADDR       string
# D2B_BOOT_IPV4         string
# D2B_BOOT_NETMASK      string
# D2B_BOOT_GATEWAY      string
# D2B_BOOT_SUBNET       string
# D2B_BUILD             flag
# D2B_HOSTNAME          string
# D2B_SITE              string ?
# D2B_DOMAIN            string
# D2B_DNS_SERVERS       list
# D2B_DHCP              flag
# D2B_IPV4              string
# D2B_NETMASK           string
# D2B_GATEWAY           string
# D2B_IMAGE             string
# D2B_IMAGE_VERSION     string
# D2B_IMAGE_TYPE        string
# D2B_GPT               string
# D2B_RAID              string
# D2B_BTRF              string

# Set key in resource
set_key() {
    local resource="$1" key="$2" value="$3" json

    json="[ { \"op\": \"replace\", \"path\": \"${key/.//}\", \"value\": \"${value}\" } ]"
    curl -s -H "Content-Type: application/json" -X PATCH -d "${json}" "${D2B_API_URL}${resource}${key/.//}"
}

set -eu

# Defaults
BASE="/root"
LOG="${BASE}/install.log"
KOPTS_DEFAULT="vconsole.keymap=us verbose nomodeset crashkernel=auto selinux=0 pcie_aspm=off"
DOCKER_DIR="/mnt/docker"

# Source functions
source "${BASE}/functions.sh"

# Setup logging
info "Logging to: ${LOG}"
exec &> >(tee -a ${LOG})

# Get options passed by PXE menu
DEBUG=$(get_kopt_flag "D2B_DEBUG")
info "Debug: $(print_bool ${DEBUG})"
[ ${DEBUG} -eq ${TRUE} ] && set -x

# Source install functions
info "Source install functions"
source "${BASE}/install_functions.sh"
source "${BASE}/install_functions-${DISTRO}.sh"

# Get options passed by PXE menu
BOOT_HWADDR=$(get_kopt "D2B_BOOT_HWADDR")
[ -n "${BOOT_HWADDR}" ] || fatal "Missing kernel option: D2B_BOOT_HWADDR"
info "Boot Hardware Address: ${BOOT_HWADDR}"

BOOT_IF=$(grep -i --with-filename ${BOOT_HWADDR} /sys/class/net/*/address | cut -d/ -f5)
[ -n "${BOOT_IF}" ] || fatal "Can't determine boot interface"
info "Boot Interface: ${BOOT_IF}"

BOOT_IPV4=$(get_kopt "D2B_BOOT_IPV4")
[ -n "${BOOT_IPV4}" ] || fatal "Missing kernel option: D2B_BOOT_IPV4"
info "Boot IPv4 Address: ${BOOT_IPV4}"

BOOT_NETMASK=$(get_kopt "D2B_BOOT_NETMASK")
[ -n "${BOOT_NETMASK}" ] || fatal "Missing kernel option: D2B_BOOT_NETMASK"
info "Boot Netmask: ${BOOT_NETMASK}"

BOOT_GATEWAY=$(get_kopt "D2B_BOOT_GATEWAY")
[ -n "${BOOT_GATEWAY}" ] || fatal "Missing kernel option: D2B_BOOT_GATEWAY"
info "Boot Gateway: ${BOOT_GATEWAY}"

BOOT_SUBNET=$(get_kopt "D2B_BOOT_SUBNET")
[ -n "$BOOT_SUBNET}" ] || fatal "Missing kernel option: D2B_BOOT_SUBNET"
info "Boot Subnet: ${BOOT_SUBNET}"

BUILD=$(get_kopt_flag "D2B_BUILD")
[ "${BUILD}" != "${TRUE}" ] && fatal "Build flag is not set, won't (re-)build this host"

HOSTNAME=$(get_kopt "D2B_HOSTNAME")
[ -n "${HOSTNAME}" ] || fatal "Missing kernel option: D2B_HOSTNAME"
info "Hostname: ${HOSTNAME}"

SITE=$(get_kopt "D2B_SITE")
if [ -z "${SITE}" ]; then
    warning "Missing kernel option: D2B_SITE"
    info "Will determine domain and DNS servers using DHCP"

    DOMAIN=$(awk '/domain/ {print $2}' /etc/resolv.conf | head -1)
    [ -n "${DOMAIN}" ] || fatal "Can't determine DHCP domain"

    DNS_SERVERS=$(awk '/nameserver/ {print $2}' /etc/resolv.conf | tr '\n' ' ')
    [ -n "${DNS_SERVERS}" ] || fatal "Can't determine DHCP nameserver(s)"
else
    info "Site: ${SITE}"

    DOMAIN=$(get_kopt "D2B_DOMAIN")
    [ -n "${DOMAIN}" ] || fatal "Missing kernel option: D2B_DOMAIN"

    DNS_SERVERS=$(get_kopt "D2B_DNS_SERVERS")
    [ -n "${DNS_SERVERS}" ] || fatal "Missing kernel option: D2B_DNS_SERVERS"
fi

info "Domain: ${DOMAIN}"
info "DNS Servers: ${DNS_SERVERS}"

DHCP=$(get_kopt_flag "D2B_DHCP")

if [ ${DHCP} == ${TRUE} ]; then
    IPV4=$(get_kopt "D2B_IPV4")
    if [ -n "${IPV4}" ]; then
        warning "Missing kernel option: D2B_IPV4"
        info "Defaulting to DHCP, since network configuration is missing"
        DHCP=${TRUE}
    else
	info "IPv4 Address: ${IPV4}"

        NETMASK=$(get_kopt "D2B_NETMASK")
        [ -n "${NETMASK}" ] || fatal "Missing kernel option: D2B_NETMASK"

        GATEWAY=$(get_kopt "D2B_GATEWAY")
        [ -n "${GATEWAY}" ] || fatal "Missing kernel option: D2B_GATEWAY"
    fi
fi

info "DHCP: $(print_bool ${DHCP})"

IMAGE=$(get_kopt "D2B_IMAGE")
[ -n "${IMAGE}" ] || fatal "Missing kernel option: D2B_IMAGE"
info "Image: ${IMAGE}"

IMAGE_VERSION=$(get_kopt "D2B_IMAGE_VERSION")
[ -n "${IMAGE_VERSION}" ] || fatal "Missing kernel option: D2B_IMAGE_VERSION"
info "Image Version: ${IMAGE_VERSION}"

IMAGE_TYPE=$(get_kopt "D2B_IMAGE_TYPE")
[ -n "${IMAGE_TYPE}" ] || fatal "Missing kernel option: D2B_IMAGE_TYPE"
info "Image Type: ${IMAGE_TYPE}"

GPT=$(get_kopt_flag "D2B_GPT")
info "GPT: $(print_bool ${GPT})"

RAID=$(get_kopt_flag "D2B_RAID")
info "RAID: $(print_bool ${RAID})"

BTRFS=$(get_kopt_flag "D2B_BTRFS")
info "BTRFS: $(print_bool ${BTRFS})"

# Set hostname
FQDN="${HOSTNAME}.${DOMAIN}"
info "Set hostname: ${FQDN}"
hostname "${FQDN}"

# Set initial root password to "dock2box"
info "Set root password"
/usr/sbin/usermod -p '$6$WyzcP9uG$64NvhY1P.W1d08CRxFjJ5QUDyZcrzyjxVPV82bAxBgf9eE2sdSZs.v47HGdWPvLxhNxAkrJten86R2Qb1vdhe/' root

# Maximize root filesystem size to appease docker IFF we have >8GB of RAM
MEM_SIZE=$(free -m | grep Mem | awk '{print $2}')

# Strip off last the chars (just show GB)
MEM_GB=$(echo "${MEM_SIZE}" | sed "s/...$//")

if [ ${MEM_GB} -gt 6 ]; then
    TGT_SIZE=$((${MEM_GB} - 2))
    echo "resize / to appropriate size"
    echo "mount -o remount size=${TGT_SIZE}G /"
    mount -o "remount,size=${TGT_SIZE}G" /
fi

# Identify disk name(s)
DISK1=$(get_sda_device)
if [[ "${SDA}" == '/dev/vda' ]]; then
   DISK2='/dev/vdb'
else
   DISK2='/dev/sdb'
fi
ROOT='/mnt/sysimage'

# Wait for first disk
info "Wait for first disk ${DISK1} to become available"
wait_for_disk ${DISK1}

# Wipe disk(s)
info "Wipe disk(s)"
wipe_disks ${DISK1} ${DISK2}

# Partition disk(s)
info 'Partition disk(s)'
if [ ${RAID} == ${TRUE} ]; then
    if [ ${GPT} == ${TRUE} ]; then
        partition_disk_gpt ${SDA}
        partition_disk_gpt ${SDB}
    else
        partition_disk_mbr ${SDA}
        partition_disk_mbr ${SDB}
    fi
    sleep 5s
    info "Create RAID 1 mirroring"
    if [ ${BTRFS} == ${TRUE} ]; then
        create_btrfs_fs ${SDA} ${SDB} raid1
    else
        create_mdraid ${SDA} ${SDB}
    fi
else
    if [ ${GPT} == ${TRUE} ]; then
        partition_disk_gpt ${SDA}
    else
        partition_disk_mbr ${SDA}
    fi
    if [ ${BTRFS} == ${TRUE} ]; then
        create_btrfs_fs ${SDA} ${SDB} no
    else
        create_lvm_vgs ${SDA} ${SDB}
    fi
fi

# Create logical volumes/sub-volumes/file-systems
if [ "${BTRFS}" == ${TRUE} ]; then
    info'Setup sub-volumes'
    create_btrfs_subvol
    info 'Mount sub-volumes'
    mount_btrfs_fss
else
    info 'Setup logical volumes'
    create_lvm_lvs
    info 'Setup file-systems'
    create_lvm_fss
    info 'Mount file-systems'
    mount_lvm_fss
fi

# Download/pull host image and apply it
if [ "${IMAGE_TYPE}" == "docker" ]; then
    info 'Create temporary docker directory'
    mkdir -p ${DOCKER_DIR}
    info 'Start Docker Engine'
    docker_start ${DOCKER_DIR}/docker.tmp ${DOCKER_REGISTRY}
    info 'Pull Docker image and apply it to disk(s)'
    docker_image_pull_and_apply ${IMAGE}:${IMAGE_VERSION} ${DOCKER_REGISTRY}
else
    info "Download image file"
    download "${IMGSTORE}/image.tar.xz" "${FILE_REGISTRY}/${IMAGE}-${IMAGE_VERSION}/image.tar.xz"
    info "Unpack system image to disk"

    # Symlinks may fail, so allow things to continue
    set +e
    tar -Jxf "${IMGSTORE}/image.tar.xz" --exclude="boot/grub/grubenv" --exclude="boot/grub2/grubenv" --exclude="sys/.keep" --warning="no-timestamp" -C "${ROOT}" --checkpoint=.5000
    set -e
fi

info "Write network configuration"
write_net_config_common
if [ ${DHCP} == ${TRUE} ]; then
    write_net_config_static
else
    write_net_config_dhcp
fi

if [ ${BTRFS} == ${TRUE} ]; then
    info "Write BTRFS fstab"
    write_fstab_btrfs
    btrfs_fix_fs
    info "Configure Docker to use BTRFS storage driver"
    configure_docker btrfs
else
    info "Write LVM fstab"
    write_fstab_lvm
    lvm_fix_fs
    info "Configure Docker to use LVM storage driver"
    configure_docker lvm
fi

# Install bootloader
info 'Install bootloader'
grub_config_and_install "${KOPTS_ALL}"

# Save logs
info 'Save logs'
save_logs

# If we're using BTRFS, take a snapshot of all subvolumes of the pristine system
if [ ${BTRFS} == ${TRUE} ]; then
    btrfs_snapshot creation
fi

# Set value to false so subsequent, accidental PXE boots don't accidentally rebuild a host
set_key "/hosts/${FQDN}" ".build" "false"

if [ ${IMAGE_TYPE} == "docker" ]; then
    docker_cleanup
fi

# dockervol is mounted regardless of image type, so cleanup
cleanup_dockervol

if [ ${DEBUG} == ${TRUE} ]; then
    warning 'Installation finished. Debug flag was set, will not reboot automatically.'
    exit 0
fi

info 'Release DHCP IP address'
if_dhcp_then_release

info 'Un-mount file-systems and reboot in 3s'
umount_fss
sleep 3
reboot
