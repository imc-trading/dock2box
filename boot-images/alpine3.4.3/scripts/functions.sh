# Get disk size in MB
disk_size_mb() {
    local disk=$1

    parted --script $disk unit MB print | grep "^Disk $disk" | awk '{print $NF}'
}

# Wait for disks to become available
wait_for_disk() {
    local disk=$1

    while [ ! -b $disk ]; do
        sleep 1
        echo -n .
    done
}

# Wait for network to become available
wait_for_network() {
    local gw=$(route -n | awk '/UG/ { print $2 }')

    counter=0
    while ! ping -w 1 -c 1 $gw &>/dev/null; do
        sleep 1s
        echo -n .
        counter=$((counter+1))
        if [[ "$counter" -gt 100 ]]; then
            fatal "Network failed to come up after 100 seconds"
        fi
    done
    echo
}

# Zero out partition table completely
clear_disk_pt() {
    local disk=$1

    dd if=/dev/zero of=$disk bs=512 count=1 conv=notrunc
}

# Wipe LVM
wipe_lvm() {
    set +e
    udevadm trigger
    udevadm settle
    vgscan
    lvscan

    for VG in $(vgs -o name | tail -n +2); do
        vgremove --force ${VG}
    done
    for PV in $(pvs -o name | tail -n +2); do
        pvremove --force --force ${PV}
    done
    sleep 1s
    set -e
}

# Partition disk using GPT
partition_disk_gpt() {
    local disk=$1

    # GUID Partition Table (GPT)
    sgdisk -Z $disk
    sgdisk -og $disk

    # Legacy BIOS boot partition is always EF02
    sgdisk -n 1:2048:+1M -t 1:ef02 $disk

    # 2nd (/boot) partition is always EF00
    sgdisk -n 2::+500M -t 2:ef00 $disk

    # LVM partition
    sgdisk -n 3 -t 3:8300 $disk

    # IMPORTANT: Ensure bios that can't handle GPT directly can still boot (e.g. Dell Desktops)
    parted --script --align optimal $disk disk_toggle pmbr_boot
}

# Partition disk using MBR
partition_disk_mbr() {
    local disk=$1

    # MBR Partition Table
    parted --script --align optimal $disk mktable msdos

    # 1st (placehoder) is always oddball (0="empty") type so we don't use it
    # this corresponds to the "Legacy BIOS" GPT partition, and is created so
    # we use the same partitioning scheme regardless of partition table type
    parted --script --align optimal $disk mkpart primary ext2 1 2
    sfdisk --part-type $disk 1 df

    # 2nd (/boot) partition is always EFI type, in case we want to boot UEFI
    parted --script --align optimal $disk mkpart primary 2 502
    parted --script $disk set 2 boot on
    sfdisk --part-type $disk 2 ef

    # LVM partition
    parted --script --align optimal $disk mkpart primary 502 $(disk_size_mb $disk)
}

# Create base layout for software RAID
create_mdraid() {
    local disk1=$1 disk2=$2
    local md0="/dev/md0" md1="/dev/md1"
    local devroot=${md1}

    # Set RAID flag
    parted --script $disk1 set 2 raid on
    parted --script $disk1 set 3 raid on
    parted --script $disk2 set 2 raid on
    parted --script $disk2 set 3 raid on

    # Create RAID
    mdadm -Ss
    yes | mdadm --create ${md0} --level=1 --raid-devices=2 ${disk1}2 ${disk2}2
    yes | mdadm --create ${md1} --level=1 --raid-devices=2 ${disk1}3 ${disk2}3

    # FAT filesystem label can only contain uppercase chars
    mkfs.vfat -F32 -n BOOTFS ${md0}
    sleep 0.2
    partprobe ${disk1} ${disk2}

    # Create the root VG
    pvcreate --force --force --yes ${devroot}
    vgcreate vg_root ${devroot}

    DEV_BOOT=${md0}
}

create_lvm_vgs() {
    local disk=$1 devroot=${1}3

    mkfs.vfat -F32 -n BOOTFS ${disk}2
    sleep 0.2
    partprobe ${disk}

    pvcreate --force --force --yes ${devroot}
    vgcreate vg_root ${devroot}

    DEV_BOOT="${disk}2"
}

# Create Logical Volumes
create_lvm_lvs() {
    [ -n "${KOPT_root_size:-}" ] || KOPT_root_size="10G"
    [ -n "${KOPT_swap_size:-}" ] || KOPT_swap_size="4G"
    [ -n "${KOPT_var_size:-}" ] || KOPT_var_size="15G"

    lvcreate -y -n lv_root -L $KOPT_root_size vg_root
    lvcreate -y -n lv_swap -L $KOPT_swap_size vg_root
    lvcreate -y -n lv_var -L $KOPT_var_size vg_root
    lvcreate -y -n lv_docker -l +100%FREE vg_root
}

# Create filesystems
create_lvm_fss() {
    mkfs.ext4 -i 4096 -q -L "rootfs" /dev/vg_root/lv_root
    mkfs.ext4 -i 4096 -q -L "varfs" /dev/vg_root/lv_var
    mkswap -L swapfs /dev/vg_root/lv_swap
    mkfs.ext4 -i 4096 -q -L "docker"  /dev/vg_root/lv_docker
}

# Mount filesystmes
mount_fs() {
    mkdir -p $SYSROOT
    mount -t ext4 /dev/vg_root/lv_root $SYSROOT

    mkdir -p $SYSROOT/boot
    mount -t vfat ${DEV_BOOT} $SYSROOT/boot

    mkdir -p $SYSROOT/var
    mount -t ext4 /dev/vg_root/lv_var $SYSROOT/var

    mkdir -p $TMPDIR
    mount -t ext4 /dev/vg_root/lv_docker $TMPDIR

    mkdir -p $SYSROOT/dev
    mount -o bind /dev $SYSROOT/dev

    mkdir -p $SYSROOT/proc
    mount -o bind /proc $SYSROOT/proc

    mkdir -p $SYSROOT/sys
    mount -o bind /sys $SYSROOT/sys
}

umount_fs() {
    umount $SYSROOT/dev
    umount $SYSROOT/proc
    umount $SYSROOT/sys
    umount $SYSROOT/boot
    umount $SYSROOT/var
    umount $SYSROOT
    umount $TMPDIR
    sync
}

add_sshkey() {
    if [ "${KOPT_sshkey:-none}" == "none" ]; then
        warn "No SSH Key specified"
        return
    fi

    mkdir -p $SYSROOT/home/dock2box/.ssh
    chmod 700 $SYSROOT/home/dock2box $SYSROOT/home/dock2box/.ssh
    echo "ssh-rsa $KOPT_sshkey sshkey@kopt" >>$SYSROOT/home/dock2box/.ssh/authorized_keys
    chroot $SYSROOT chown -R dock2box:dock2box /home/dock2box
}

get_sda_device() {
    sdadev=$(ls /dev/vda)
    if [[ -z "${sdadev}+abc" ]] || ( [[ -z "${sdadev}" ]] && [[ "${sdadev+abc}" == "abc" ]] ); then
        sdadev='/dev/sda'
    fi
    echo ${sdadev}
}

# Add error handling for HTTP status

docker_download() {
    local reg=$1 name=$2 target=$3 parallel=$4 header="${5}" full short checksum csize
    shift; shift; shift; shift; shift

    while [ -n "$1" ]; do
        full=${1##sha256:}
        short=${full:0:12}

        file="$target/${1}.tar.gz"
        if [ -f "$file" ]; then
            checksum=$( sha256sum $file | awk '{ print $1 }' )
            if [ "$checksum" != "$full" ]; then
                warn "Re-download layer $short" 10
                if [ -e /tmp/progress ]; then
                    csize=$( curl -sI "https://${reg}/v2/${name}/blobs/${1}" | awk '/Content-Length/ {print $2}' )
                    curl -s -L -H "$header" "https://${reg}/v2/${name}/blobs/${1}" | pv -n -s $csize >$file 2>/tmp/progress
                else
                    curl --progress-bar -o $file -L -H "$header" "https://${reg}/v2/${name}/blobs/${1}"
                fi
            else
                info "Skip layer $short"
            fi
        else
            info "Download layer $short" 10
            if [ -e /tmp/progress ]; then
                csize=$( curl -sI "https://${reg}/v2/${name}/blobs/${1}" | awk '/Content-Length/ {print $2}' )
                curl -s -L -H "$header" "https://${reg}/v2/${name}/blobs/${1}" | pv -n -s $csize >$file 2>/tmp/progress
            else
                curl --progress-bar -o $file -L -H "$header" "https://${reg}/v2/${name}/blobs/${1}"
            fi
        fi
        shift
    done
}

docker_apply() {
    local source=$1 full short
    shift

# Exclude should be in image creation to make it generic
# Also fix /var/lib/yum/yumdb

    while [ -n "$1" ]; do
        full=${1##sha256:}
        short=${full:0:12}

        info "Apply layer ${short}" 0
        file="$source/${1}.tar.gz"
        if [ -e /tmp/progress ]; then
            ( pv -n $file | tar xzf - -h -C $SYSROOT \
            --exclude=./boot/grub/grubenv \
            --exclude=./boot/grub2/grubenv \
            --exclude=./sys \
            --exclude=./dev \
            --exclude=./proc
            ) >/tmp/progress 2>&1
        else
            tar -xzf $file -h -C $SYSROOT \
            --exclude=./boot/grub/grubenv \
            --exclude=./boot/grub2/grubenv \
            --exclude=./sys \
            --exclude=./dev \
            --exclude=./proc
        fi
        shift
    done
}

docker_pull() {
    local reg=$1 name=$2 tag=$3 target=$4 parallel=$5 header="Accept: application/json"

    # Get auth token if we're using public registry
    if [ "$reg" == "registry.hub.docker.com" ]; then
        # Get auth token
        token=$( curl -s "https://auth.docker.io/token?service=registry.docker.io&scope=repository:${name}:pull" | jq -r .token )
        header="Authorization: Bearer $token"
    fi

    # Get layers
    manifest=$( curl -s -H "$header" "https://${reg}/v2/${name}/manifests/${tag}" )
    schemaVer=$( echo $manifest | jq -r .schemaVersion )

    case $schemaVer in
        1) layers=$(curl -s -H "$header" "https://${reg}/v2/${name}/manifests/${tag}" | jq -r .fsLayers[].blobSum ) ;;
        2) layers=$(curl -s -H "$header" "https://${reg}/v2/${name}/manifests/${tag}" | jq -r .layers[].digest ) ;;
        *) fatal "Unsupported schema version: $schemaVer" ;;
    esac

    docker_download $reg $name $target $parallel "${header}" $layers

    # Run twice in-case of failure
    docker_download $reg $name $target $parallel "${header}" $layers

    docker_apply $target $layers
}
