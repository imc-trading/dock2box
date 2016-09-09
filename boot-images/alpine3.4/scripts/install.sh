#!/bin/bash

set -eu

LOG="./install.log"
LCK="./.dock2box"

# Get kernel options
set -- $(cat /proc/cmdline)
myopts="debug nocolor distro image registry ip gpt mdraid install sshkey"

for opt; do
    for i in $myopts; do
        case "$opt" in
            $i=*)   eval "KOPT_${i}='${opt#*=}'";;
            $i)     eval "KOPT_${i}=true";;
        esac
    done
done

# Include functions
source "./functions.sh"

if [ -e $LCK ]; then
    fatal "Installation already running or completed, to run again remove lockfile: $LCK"
fi
touch $LCK

if [ "${KOPT_install:-false}" == "false" ]; then
    fatal "Kernel option \"install\" isn't set exiting"
fi

# Include distribution functions
source "./functions-${KOPT_distro}.sh"

# Setup logging
info "Logging to: $LOG"
exec &> >(tee -a $LOG)

# Enable debug
if [ "${KOPT_debug:-false}" == "true" ]; then
    set -x
fi

# Get first disk
disk0=$(get_sda_device)
if [ "$disk0" == '/dev/vda' ]; then
   disk1='/dev/vdb'
else
   disk1='/dev/sdb'
fi

# Get first interface
intf="eth0"

# Get hardware address
hwaddr=$(cat /sys/class/net/$intf/address)

info 'Wait for disk to become available'
wait_for_disk $disk0

info 'Stop Docker and unmount partitions in-case we''re re-running this manually'
set +e
stop_docker
umount_fs
set -e

info 'Wipe disk(s)'
wipe_lvm
clear_disk_pt $disk0

if [ "${KOPT_mdraid:-false}" == "true" ]; then
    wipe_mdraid $disk0 $disk1
    clear_disk_pt $disk1
fi

info 'Partitions disk(s)'
if [ "${KOPT_mdraid:-false}" == "true" ]; then
    info 'Using mdraid for mirroring'
    if [ "${KOPT_gpt:-false}" == "true" ]; then
        info 'Using GPT'
        partition_disk_gpt $disk0
        partition_disk_gpt $disk1
    else
        partition_disk_mbr $disk0
        partition_disk_mbr $disk1
    fi
    sleep 5s
    create_mdraid $disk0 $disk1
else
    if [ "${KOPT_gpt:-false}" == "true" ]; then
        info 'Using GPT'
        partition_disk_gpt $disk0
    else
        partition_disk_mbr $disk0
    fi
    create_lvm_vgs $disk0
fi

info 'Setup logical volumes'
create_lvm_lvs

info 'Setup filesystems'
create_lvm_fss

info 'Mount filesystems'
mount_fs

info 'Start Docker'
start_docker

# Some network cards reset the link when the Docker bridge interface is added
sleep 5
wait_for_network

info 'Pull Docker image and apply to disk'
apply_docker_image $KOPT_image $KOPT_registry

info 'Write network config'
conf_net
if [ "$KOPT_ip" == "dhcp" ]; then
    conf_net_dhcp $intf $hwaddr
#else
#    write_net_config_static
fi

info "Add SSH Key"
add_sshkey

info "Write fstab"
write_fstab

info 'Install bootloader'
install_grub $disk0 "$KOPTS ${KOPT_kopts:-}"

stop_docker
umount_fs

if [ "${KOPT_debug:-false}" == "true" ]; then
    info "Debug set skipping reboot"
    while true; do
        info "Stop service dock2box to exit script or press [CTRL+C] when running it manually."
        sleep 5
    done
fi

sleep 3
release_dhcp
reboot
