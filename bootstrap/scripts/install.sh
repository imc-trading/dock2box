#!/bin/bash

set -eu

info() {
    local message="$1"

    echo "${message}" >&2
}

fatal() {
    local message="$1"

    echo "^[[1;31m${message}^[[0m" >&2
    sleep 60
    exit 1
}

download() {
    local url="$1" target="$2"

    wget --quiet --dns-timeout=2 --connect-timeout=2 --read-timeout=2 --no-check-certificate --output-document "${target}" "${url}"
}

set_key() {
    local resource="$1" key="$2" value="$3"

    read -d '' json <<EOF
[
  {
    "op": "replace",
    "path": "${key/.//}",
    "value": "${value}"
  }
]
EOF

    curl -s -H "Content-Type: application/json" -X PATCH -d "${json}" " "${API_URL}${resource}${key/.//}"
}

API_URL="http://dock2box:8080/api/v1"
BASE="/root"
LOG="${BASE}/install.log"
ARGS=$*
KOPTS_DEFAULT="vconsole.keymap=us verbose nomodeset crashkernel=auto selinux=0 pcie_aspm=off"
DOCKER_DIR="/mnt/docker"

# Set debug
if [[ "${DEBUG}" == "true" ]]; then
    set -x
fi

# Setup logging
info "Logging to: ${LOG}"
exec &> >(tee -a ${LOG})

# Check that we're allowed to build this box
[[ ! "${BUILD}" =~ ^[Tt]rue$ ]] && fatal "BUILD FLAG IS NOT SET TO TRUE, WILL NOT (RE-)BUILD THIS HOST"

# Check options
[[ "${IMAGE_TYPE}" =~ ^(docker|file)$ ]] && fatal "INCORRECT IMAGE TYPE: ${IMAGE_TYPE}, WILL NOT (RE-)BUILD THIS HOST"

# Download functions
cd ${BASE}
download "${SCRIPT/install.sh/functions.sh}" "${BASE}/functions.sh"
download "${SCRIPT/install.sh/functions-${DISTRO}.sh}" "${BASE}/functions-${DISTRO}.sh"

# Source functions
source "${BASE}/functions.sh"
source "${BASE}/functions-${DISTRO}.sh"

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
if [[ "${RAID}" =~ ^[Tt]rue$ ]]; then
    if [[ "${GPT}" =~ ^[Tt]rue$ ]]; then
        info "Use GPT partititions"
        partition_disk_gpt ${SDA}
        partition_disk_gpt ${SDB}
    else
        partition_disk_mbr ${SDA}
        partition_disk_mbr ${SDB}
    fi
    sleep 5s
    info "Create software RAID-1 mirroring"
    if [[ "${VOLMGT}" == "btrfs" ]]; then
        info "Use BTRFS for volume management"
        create_btrfs_fs ${SDA} ${SDB} raid1
    else
        create_mdraid ${SDA} ${SDB}
    fi
else
    if [[ "${GPT}" =~ ^[Tt]rue$ ]]; then
        info "Use GPT partititions"
        partition_disk_gpt ${SDA}
    else
        partition_disk_mbr ${SDA}
    fi
    if [[ "${VOLMGT}" == "btrfs" ]]; then
        info "Use BTRFS for volume management"
        create_btrfs_fs ${SDA} ${SDB} no
    else
        create_lvm_vgs ${SDA} ${SDB}
    fi
fi

# Create logical volumes/sub-volumes/file-systems
if [[ "${VOLMGT}" == "btrfs" ]]; then
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

# Dowbload/pull host image and apply it
if [[ "${IMAGE_TYPE}" =~ "docker" ]]; then
    info 'Create temporary docker directory'
    mkdir -p ${DOCKER_DIR}
    info 'Start Docker Engine'
    docker_start ${DOCKER_DIR}/docker.tmp ${DOCKER_REGISTRY}
    info 'Pull Docker image and apply it to disk(s)'
    docker_image_pull_and_apply ${IMAGE}:${IMAGE_VERSION} ${DOCKER_REGISTRY}
elif [[ "${IMAGE_TYPE}" == "file" ]]; then
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
if [[ "${DHCP}" =~ ^[Ff]alse$ ]]; then
    write_net_config_static
else
    write_net_config_dhcp
fi

if [[ "${VOLMGT}" == "btrfs" ]]; then
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
if [[ "${VOLMGT}" == "btrfs" ]]; then
    btrfs_snapshot creation
fi

# Set value to false so subsequent, accidental PXE boots don't accidentally rebuild a host
set_key "/hosts/${FQDN}" ".build" "false"

if [[ "${IMAGE_TYPE}" =~ "docker" ]]; then
    docker_cleanup
fi

# dockervol is mounted regardless of image type, so cleanup
cleanup_dockervol

if [[ "${DEBUG}" == "true" ]]; then
    msg 'Installation finished. Debug flag was set, will not reboot automatically.'
else
    if_dhcp_then_release
    if grep -q "X[tT]rue" <<< "X${KEXEC-}"; then
        msg 'KEXEC enabled, load and boot the new kernel'
        kexec_boot "${KOPTS_ALL}"
    else
        msg 'KEXEC disabled, unmount filesystems and reboot in 3s'
        umount_fss
        sleep 3
        reboot
    fi
fi
