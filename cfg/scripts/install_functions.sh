# Redirect docker storage
docker_start() {
    local dir=${1}
    local registry=${2}

# Get certificates
#    mkdir -p ${dir} /etc/docker/certs.d/${registry}
#    curl -s ... >/etc/docker/certs.d/${registry}/ca.crt
#    export DOCKER_CERT_PATH=/etc/docker

    docker -g ${dir} -d &
    sleep 5

    # Some network cards reset the link when the Docker bridge interface is added
    wait_for_network
}

# Pull Docker image and run it
# Exclude gruben symlink since we use FAT32 for UEFI /boot compatability, which does not support symlinks
docker_image_pull_and_apply() {
    local img=${1} registry=${2}
    docker run -i --entrypoint /bin/true ${registry}/${img}
    local container=$(docker ps -a -l | awk '/true/ {print $1}')

    docker export ${container} | tar -x -h -C ${ROOT} --exclude=boot/grub/grubenv --exclude=boot/grub2/grubenv --exclude=sys/.keep --warning=no-timestamp --checkpoint=.5000
}

# This fixes legacy crap that gets dumped into /fs in a way incompatible with normal filesystem layouts
# This will no longer be required once the images are cleaned up
btrfs_fix_fs() {
    if [ -L ${ROOT}/usr/portage ]; then
        rm -rf ${ROOT}/usr/portage
    fi
    if [ -d ${ROOT}/fs/portage ]; then
        mv ${ROOT}/fs/portage ${ROOT}/usr
    fi
    if [ -d ${ROOT}/fs/var ]; then
        test -d ${ROOT}/fs/var/db && mv ${ROOT}/fs/var/db ${ROOT}/var
        test -l ${ROOT}/fs/var/mail && mv ${ROOT}/fs/var/mail ${ROOT}/var
        test -l ${ROOT}/fs/var/lock && mv ${ROOT}/fs/var/lock ${ROOT}/var
        test -d ${ROOT}/fs/var/spool && mv ${ROOT}/fs/var/spool ${ROOT}/var
        test -l ${ROOT}/fs/var/mail && mv ${ROOT}/fs/var/run ${ROOT}/var
        test -d ${ROOT}/fs/var/lib && cp -a ${ROOT}/fs/var/lib ${ROOT}/var
        test -d ${ROOT}/fs/var/cache && cp -a ${ROOT}/fs/var/cache ${ROOT}/var
        test -d ${ROOT}/fs/var/log && cp -a ${ROOT}/fs/var/log ${ROOT}/var
        rm -rf ${ROOT}/fs/var
    fi
    if [ -L ${ROOT}/opt ]; then
	rm -rf ${ROOT}/opt
	mkdir ${ROOT}/opt
    fi
    rm -rf ${ROOT}/fs
    
}

# Transform standard layout to /fs unified partition layout
lvm_fix_fs() {
    echo "Adjusting filesystem for unified /fs partition"
    # home
    cp -a ${ROOT}/home ${ROOT}/fs/home ; rm -rf ${ROOT}/home ; cd ${ROOT} ; ln -s fs/home home ; cd -
    # opt
    cp -a ${ROOT}/opt ${ROOT}/fs/opt ; rm -rf ${ROOT}/opt ; cd ${ROOT} ; ln -s fs/opt opt ; cd -
    # var
    cp -a ${ROOT}/var ${ROOT}/fs/var ; mv ${ROOT}/var ${ROOT}/not.var ; cd ${ROOT} ; ln -s fs/var var ; cd -
    # separate log partition already mounted.  This is painful, but too many people don't monitor or limit log sizes
    mv ${ROOT}/fs/var/log ${ROOT}/fs/var/log.old
    cd ${ROOT}/fs/var && ln -s ../log && cd -
    mv ${ROOT}/fs/var/log.old/* ${ROOT}/fs/log
    # run
    cd ${ROOT}/fs ; ln -s ../run run ; cd -
    # data
    mkdir -p ${ROOT}/fs/data
    cd ${ROOT} ; ln -s fs/data data ; cd -
}

# Cleanup after docker stage
docker_cleanup() {
    docker stop $(docker ps -a -q)
    docker rm -v $(docker ps -a -q)
    pkill -f docker
    sleep 3
}

# Extract boot args from /proc/cmdline
get_cmdline_arg() {
    local name=${1}
    local val

    if grep -w ${name} /proc/cmdline &>/dev/null; then
        # Use printf to catch multiword/quoted arg values
        val=$(xargs printf "%s\n" </proc/cmdline | grep ^${name}= | cut -d '=' -f 2-)
        if [ -n "${val}" ]; then
            echo "${val}"
        else
            echo "~"
        fi
    fi
}

# Get disk size in MB
disk_size_mb() {
    local disk=$1

    parted --script ${disk} unit MB print | grep "^Disk $disk" | awk '{print $NF}'
}

# Wait for disk to become available
wait_for_disk() {
    local disk=$1
    
    while [ ! -b $disk ]; do
        sleep 1
        echo -n .
    done
}

# Wait for network to become available
wait_for_network() {
    counter=0
    while ! ping -w 1 -c 1 ${GATEWAY} &>/dev/null; do
    # Please document what the following line does and how.  3 of us could not make head nor tail of it:
	echo "wait_for_network, ${counter}. attempt"
        sleep 1
        echo -n .
        counter=$((counter+1))
        if [[ "${counter}" -gt 100 ]]; then
            echo "The network has failed to come up after 100 seconds.  Giving up"
            exit 1
        fi
    done
    echo
}

# Get network device name associated with mac address.  Filters out any duplicate bridge (brX) devices, bondX devices, and 
# duplicate devices resulting from bonds (by ensuring the lowest device is returned, which should match where bond gets its mac address) 
get_device_from_mac() {
    local mac=$1
    LIST=$(ip addr show |grep -B 1 ${mac} | grep -v ${mac} | sed '/^--$/d' | awk '{print $2}' | cut -d : -f 1 | egrep -v '^br[0-9]' | grep -v '^bond[0-9]' | sort -u)
    echo ${LIST} | cut -d ' ' -f 1 
}

# Predict "predictable"
get_predictable_name() {
    local mac=$1

    local olddev=$(get_device_from_mac ${mac})

    # if things crap out, we'll return an empty string.  This should never happen though.
    local NEWDEV=""

    # In order of precedence, since "predicatable" device names aren't really all that predictable
    # 1. Names incorporating Firmware/BIOS provided index numbers for on-board devices (ID_NET_NAME_ONBOARD, example: eno1)
    local OUT=$(udevadm test-builtin net_id /sys/class/net/${olddev} 2> /dev/null | grep ID_NET_NAME_ONBOARD)
    local EXIT=$?
    if [ ${EXIT} -eq 0 ]; then
        NEWDEV=$(echo ${OUT} | cut -d = -f 2)
    else
        # 2. Names incorporating Firmware/BIOS provided PCI Express hotplug slot index numbers (ID_NET_NAME_SLOT, example: ens1)
        OUT=$(udevadm test-builtin net_id /sys/class/net/${olddev} 2> /dev/null | grep ID_NET_NAME_SLOT)
        EXIT=$?
        if [ ${EXIT} -eq 0 ]; then
            NEWDEV=$(echo ${OUT} | cut -d = -f 2)
        else
            # 3. Names incorporating physical/geographical location of the connector of the hardware (ID_NET_NAME_PATH, example: enp2s0)
            OUT=$(udevadm test-builtin net_id /sys/class/net/${olddev} 2> /dev/null | grep ID_NET_NAME_PATH)
            EXIT=$?
            if [ ${EXIT} -eq 0 ]; then
                NEWDEV=$(echo ${OUT} | cut -d = -f 2)
            else
                # 4. rule 4 is disabled by default (but would use mac address as device name, a la: enx78e7d1ea46da), so we fall back to ...
                # 5. Classic kernel-native ethX naming (example: eth0), also not predictable and not stable across reboots without udev rule hacks
                NEWDEV=${dev}
            fi
        fi
    fi
    echo ${NEWDEV}
}

# zero out partition table completely
clear_disk_pt() {
    local disk=$1

    dd if=/dev/zero of=$disk bs=512 count=1 conv=notrunc
}

wipe_lvm()
{
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
    sleep 1
    set -e
}

wipe_mdraid() {
    local disk1=$1
    local disk2=$2

    set +e

    echo "WARNING: This probably requires distro-specific actions, but no wipe_mdraid() function was defined for your distro!"
    echo "I'll do my best, but no guarantees!"

    # Cleanup software RAID
    for MD in $(cd /dev; ls -d md? 2>/dev/null); do
        mdadm --stop ${MD}
    done

    # clean up our devices
    udevadm control --reload
    sleep 1
    set -e
}

# Wipe anything on disk
wipe_disks() {
    local disk1=$1
    local disk2=$2

    set +e
    wipe_mdraid $disk1 $disk2
    clear_disk_pt $disk2
    clear_disk_pt $disk1
    set -e
}

# Partition disk using GPT
partition_disk_gpt() {
    local disk=$1

    # GUID Partition Table (GPT)
    #parted --script --align optimal $disk mktable gpt
    sgdisk -Z $disk
    sgdisk -og $disk

    # Legacy BIOS boot partition is always EF02
    #parted --script --align optimal $disk mkpart primary ext2 1 2
    sgdisk -n 1:2048:+1M -t 1:ef02 $disk
    #sgdisk $disk --typecode=1:EF02
    lastsector=$(gdisk -l $disk | grep EF02 | awk '{print $3}')

    # 2nd (/boot) partition is always EF00
    parted --script --align optimal $disk mkpart primary fat32 2 502
    sgdisk -n 2:${lastsector}:+500M -t 2:ef00
    #sgdisk $disk --typecode=2:EF00
    lastsector=$(gdisk -l $disk | grep EF00 | awk '{print $3}')

    # LVM partition
    #parted --script --align optimal $disk mkpart primary 502 $(disk_size_mb $disk)
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

    # 2nd (/boot) artition is always EFI type, in case we want to boot UEFI
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

    # Set up some global vars
    DEV_BOOT=${md0}
    VG_FS="vg_root"
    LV_FS_SIZE="-l 80%FREE"
}

# helper function to mount special filesystems required for chroot to work correctly
chroot_mounts() {
    local btrfs=$1

    # Chroot dependencies
    mkdir -p ${ROOT}/{dev,proc,sys}
    mount -o bind /dev ${ROOT}/dev
    mount -o bind /proc ${ROOT}/proc
    mount -o bind /sys ${ROOT}/sys
    # docker dirs
    mkdir -p /mnt/docker
    if [ ${btrfs} == ${TRUE} ]; then
        btrfs subvolume create /btrfs/${DATALABEL}/tmpdocker
        chattr -R +C /btrfs/${DATALABEL}/tmpdocker
	mount -L ${DATALABEL} -t btrfs -o defaults,noatime,autodefrag,nodatacow,subvol=tmpdocker /mnt/docker
    else
    	mount -L tmpdocker /mnt/docker
    fi
}

# prepare mountpoint and snapshots
prepare_btrfs() {
    label=$1
    mnt=/btrfs/${label}
    mkdir -p ${mnt}
    mount -L ${label} ${mnt}

    mkdir ${mnt}/snapshots
}

# Create BTRFS filesystem
create_btrfs_fs() {
    local disk1=$1 disk2=$2 
    BTRFSRAID=$3
    BTRFSDEV1=${disk1}3
    BTRFSDEV2=${disk2}3

    if [ -b ${disk2} ]; then
        # if raid, then assume 2nd disk is partitioned like first disk (this covers our use case but
        # should really be refactored)
        if [[ "${BTRFSRAID}" == "raid1" ]]; then
            mkfs.btrfs -f -L d2n -m ${BTRFSRAID} -d ${BTRFSRAID} ${BTRFSDEV1} ${BTRFSDEV2}
        # if not raid, then add a 2nd btrfs for second drive
        else
            mkfs.btrfs -f -L d2b ${disk1}3
            mkfs.btrfs -f -L d2b-drive2 ${disk2}
            # add 2nd mount here
            prepare_btrfs d2b-drive2
        fi
    else
        mkfs.btrfs -f -L d2b ${disk1}3
    fi
    prepare_btrfs d2b
}

# Create BTRFS subvolumes
# TODO: preserve rebuilds by detecting root0, etc. and incrementing to root1 etc
create_btrfs_subvol() {
    DATALABEL="d2b2" # deliberately wrong to detect bugs
    BTRFSROOT=root0
    btrfs subvolume create /btrfs/d2b/${BTRFSROOT}
    btrfs quota enable /btrfs/d2b/${BTRFSROOT}
    btrfs subvolume create /btrfs/d2b/home0
    btrfs quota enable /btrfs/d2b/home0
    btrfs subvolume create /btrfs/d2b/log0
    btrfs quota enable /btrfs/d2b/log0
    if [ -d /btrfs/d2b-drive2 ]; then
        DATALABEL="d2b-drive2"
        btrfs subvolume create /btrfs/${DATALABEL}/data0
        btrfs quota enable /btrfs/${DATALABEL}/data0
        chattr -R +C /btrfs/${DATALABEL}/data0
    else
        DATALABEL="d2b"
        btrfs subvolume create /btrfs/d2b/data0
        btrfs quota enable /btrfs/d2b/data0
        chattr -R +C /btrfs/d2b/data0
    fi
}

global_mountpoints() {
    BOOT=${ROOT}/boot
    VAR=${ROOT}/var
    LOG=${ROOT}/var/log
    FSHOME=${ROOT}/home
    DATA=${ROOT}/data
    OPT=${ROOT}/opt
    IMGSTORE=/mnt/docker
}

mount_btrfs_fss() {
    global_mountpoints

    # /
    mkdir -p ${ROOT}

    # if raid != 'no', use devices
    if [[ "${BTRFSRAID}" == "raid1" ]]; then
        mount -L d2b -t btrfs -o defaults,noatime,space_cache,compress=lzo,subvol=${BTRFSROOT},device=${BTRFSDEV1},device=${BTRFSDEV2} ${ROOT}
    else
        mount -L d2b -t btrfs -o defaults,noatime,space_cache,compress=lzo,subvol=${BTRFSROOT} ${ROOT}
    fi
    # UEFI will need this
    # /boot
    # mkdir -p ${BOOT}
    # mount ${DEV_BOOT} ${BOOT}
    mkdir -p ${FSHOME}
    mount -L d2b -t btrfs -o defaults,noatime,space_cache,compress=lzo,subvol=home0 ${FSHOME}
    # /var/log
    mkdir -p ${LOG}
    mount -L d2b -t btrfs -o defaults,noatime,space_cache,compress=lzo,subvol=log0 ${LOG}
    # /data
    mkdir -p ${DATA}
    mount -L ${DATALABEL} -t btrfs -o defaults,noatime,autodefrag,nodatacow,subvol=data0 ${DATA}

    # thwart /fs/* links, since we don't need those for btrfs
    mkdir -p ${OPT}

    # temp image location
    mkdir -p ${IMGSTORE} 
    chroot_mounts
}

# Create base layout for normal servers
create_lvm_vgs() {
    local disk1=$1 disk2=$2
    local devroot=${disk1}3

    mkfs.vfat -F32 -n BOOTFS ${disk1}2
    sleep 0.2
    partprobe ${disk1}

    pvcreate --force --force --yes ${devroot}
    vgcreate vg_root ${devroot}

    # Set up some global vars
    DEV_BOOT="${disk1}2"
    if [ -b ${disk2} ]; then
        VG_FS="vg_fs"
        LV_FS_SIZE="-l 100%FREE"
        pvcreate --force --force --yes ${disk2}
        vgcreate ${VG_FS} ${disk2}
    else
        VG_FS="vg_root"
        LV_FS_SIZE="-l 80%FREE"
    fi
}

# Create Logical Volumes
create_lvm_lvs() {
    lvcreate -y -n lv_root -L 10G vg_root
    lvcreate -y -n lv_swap -L 4G vg_root
    lvcreate -y -n lv_log -L 4G vg_root
    lvcreate -y -n lv_fs ${LV_FS_SIZE} ${VG_FS}
    lvcreate -y -n lv_tmpdocker -l +100%FREE vg_root
}

# Create filesystems
# increase inodes as a precaution to allow for filesystem growth
# we should consider using xfs for /fs
create_lvm_fss() {
    mkfs.ext4 -i 4096 -q -L "rootfs" /dev/vg_root/lv_root
    mkfs.ext4 -i 4096 -q -L "log" /dev/vg_root/lv_log
    mkfs.ext4 -i 4096 -q -L "fs" /dev/${VG_FS}/lv_fs
    mkswap -L swapfs /dev/vg_root/lv_swap
    mkfs.ext4 -i 4096 -q -L "tmpdocker"  /dev/vg_root/lv_tmpdocker
}

# Mount filesystmes under ${ROOT}
# current scheme:
# /
# /fs
# /fs/log
# /opt -> /fs/opt
# /home -> /fs/home
# /data -> /fs/data
mount_lvm_fss() {

    global_mountpoints

    # Overriding global_mountpoints for /fs scheme
    VAR=${ROOT}/fs/var
    VARLOG=${ROOT}/fs/log
    FS=${ROOT}/fs
    FSHOME=${FS}/home
    OPT=${FS}/opt
    DATA=${FS}/data
    # /
    mkdir -p ${ROOT}
    mount -L rootfs ${ROOT}
    # /boot
    mkdir -p ${BOOT}
    mount ${DEV_BOOT} ${BOOT}
    # /fs
    mkdir -p ${FS}
    mount -L fs ${FS}
    # the following is handled by lvm_fix_fs():
    #mkdir -p ${FSHOME} ${OPT} ${DATA} ${IMGSTORE} ${VAR} ${VARLOG}
    ## /var/log -> /fs/log
    #(cd ${ROOT}; ln -sf fs/var var)
    #(cd ${FS}; ln -sf ../run)
    #(cd ${VAR}; ln -sf ../log)
    mkdir -p ${VARLOG}
    mount -L log ${VARLOG}
    ## replace defunct dirs with symlinks to comptability directories
    #rm -rf ${ROOT}/opt ${ROOT}/home ${ROOT}/data
    #(cd ${ROOT}; ln -sf fs/data data)
    #(cd ${ROOT}; ln -sf fs/home home)
    #(cd ${ROOT}; ln -sf fs/opt opt)

    chroot_mounts
}

# Unmount filesystems under ${ROOT}
umount_fss_lvm() {
    docker_cleanup
    umount ${ROOT}/dev ${ROOT}/proc ${ROOT}/sys ${ROOT}/boot ${ROOT}/fs/log ${ROOT}/fs ${ROOT}
    sync
}

cleanup_dockerlvm() {
    umount /mnt/docker
    sleep 3
    lvchange -an /dev/mapper/vg_root-lv_tmpdocker
    lvremove -y /dev/mapper/vg_root-lv_tmpdocker
}

cleanup_dockerbtrfs() {
    umount /mnt/docker
    btrfs subvolume delete /btrfs/${DATALABEL}/tmpdocker
}

cleanup_dockervol() {
    if [[ "${VOLMGT}" == "btrfs" ]]; then
	cleanup_dockerbtrfs
    else
	cleanup_dockerlvm
    fi
}

####  CONTINUE HERE ####

# Inject /etc/hosts
write_net_config_common() {
    echo ${HOSTNAME} >${ROOT}/etc/hostname
    cat <<EOF >${ROOT}/etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
EOF
    rm -f ${ROOT}/etc/sysconfig/network-scripts/{ifcfg-,routes}*
}

# Inject net config (static)
write_net_config_static() {
    echo "THIS IS VERY DISTRO SPECIFIC"
    echo "no write_net_config_static() function defined for your distro!"
    echo "No network configured!"
}

# Inject net config (dhcp)
write_net_config_dhcp() {
    echo "THIS IS VERY DISTRO SPECIFIC"
    echo "No write_net_config_dhcp() function defined for your distro!"
    echo "No DHCP network configuration created!"
}

# Release lease (if using DHCP)
if_dhcp_then_release() {
    if [ -z "${IP}" ]; then
        DHCLIENT=$(which dhclient)
        if [[ "${DHCLIENT}" != '' ]]; then
            dhclient -r
        fi
    fi
}

# get first disk device
get_sda_device() {
    sdadev=$(ls /dev/vda)
    if [[ "${sdadev}" == '' ]]; then
        sdadev='/dev/sda'   
    fi
    echo ${sdadev}
}

# Figure out the complete list of boot args
get_boot_args() {
    echo "THIS IS DISTRO SPECIFIC"
    echo "no get_boot_args() function defined for your distro!"
    echo "No boot arg list obtained."
}

# Inject bootloader config
grub_config_and_install() {
    echo "THIS IS DISTRO SPECIFIC"
    echo "no grub_config_and_install() function defined for your distro!"
    echo "No bootloader Installed!"
    echo "YOUR SYSTEM PROBABLY WILL NOT BOOT."
}

# take a snapshot of our btrfs subvolumes
btrfs_snapshot() {
    tag=$1
    # snapshot all our subvolumes.  This, like fstab, needs to be stored in a global variable and/or etcd/mongo entry
    # in case it doesn't already exist:
    mkdir -p /btrfs/d2b/snapshots
    btrfs sub snapshot /btrfs/d2b/${BTRFSROOT} /btrfs/d2b/snapshots/${BTRFSROOT}.${tag}
    btrfs sub snapshot /btrfs/d2b/home0 /btrfs/d2b/snapshots/home0.${tag}
    btrfs sub snapshot /btrfs/d2b/log0 /btrfs/d2b/snapshots/log0.${tag}
    # in case it doesn't already exist.  Redundant for single-drive or raid1 systems, but critical for 2-drive systems:
    mkdir -p /btrfs/${DATALABEL}/snapshots
    btrfs sub snapshot /btrfs/${DATALABEL}/data0 /btrfs/${DATALABEL}/snapshots/data0.${tag}
}

# Inject fstab - BTRFS
write_fstab_btrfs() {
    if [[ "${BTRFSRAID}" == "raid1" ]]; then
        cat << EOF > ${ROOT}/etc/fstab
LABEL=d2b       /       btrfs   defaults,noatime,space_cache,compress=lzo,subvol=${BTRFSROOT},device=${BTRFSDEV1},device=${BTRFSDEV2} 0 0
EOF
    else
        cat << EOF > ${ROOT}/etc/fstab
LABEL=d2b       /       btrfs   defaults,noatime,space_cache,compress=lzo,subvol=${BTRFSROOT} 0 0
EOF
    fi
    cat << EOF >> ${ROOT}/etc/fstab
LABEL=d2b       /home   btrfs   defaults,noatime,space_cache,compress=lzo,subvol=home0 0 0
LABEL=d2b       /var/log btrfs   defaults,noatime,space_cache,compress=lzo,subvol=log0 0 0 
# /data uses nodatacow for better performance. Recommended for mysql, VMs, lxc/docker images, etc.
LABEL=${DATALABEL} /data  btrfs   defaults,noatime,space_cache,autodefrag,nodatacow,subvol=data0 0 0
# /btrfs/* mounts hold the toplevel btrfs filesystem for snapshots, subvolumes, backups, etc
LABEL=d2b       /btrfs/d2b btrfs defaults,noatime 0 0
EOF
    mkdir -p ${ROOT}/btrfs/d2b
    if [[ "${DATALABEL}" != "d2b" ]]; then
        mkdir -p ${ROOT}/btrfs/${DATALABEL}
	echo "LABEL=${DATALABEL} /btrfs/${DATALABEL} btrfs defaults,noatime 0 0" >> ${ROOT}/etc/fstab
    fi
}

# Injext fstab - Legacy LVM
write_fstab_lvm() {
    cat <<EOF >${ROOT}/etc/fstab
LABEL=rootfs    /        ext4      defaults         1 1
LABEL=log       /fs/log  ext4      defaults         1 1
LABEL=fs        /fs      ext4      defaults         1 2
LABEL=BOOTFS    /boot    vfat      defaults         1 2
LABEL=swapfs    swap     swap      defaults         0 0
EOF
}

write_fstab() {
    if [[ "${VOLMGT}" == "btrfs" ]]; then
        msg "Write BTRFS fstab ..."
        write_fstab_btrfs
    else
        msg "Write LVM fstab ..."
        write_fstab_lvm
    fi
}

configure_docker() {
    echo "Placeholder.  This will do nothing unless you override for your distro"
}

write_build_info() {
    IDATE=$(date --rfc-3339=ns)
    cat << EOF >> ${ROOT}/etc/d2b/buildinfo.d2b

[System Build]

INSTALL_DATE: ${IDATE}
EOF
}

# Inject install logs
save_logs() {
    mkdir -p ${ROOT}${ETCD2B}/{state-when-imaged,logs,post}
    cp -a ${BASE}/*.{log,sh} ${ROOT}${ETCD2B}/state-when-imaged/
    cp /bootstrap.log ${ROOT}${ETCD2B}/logs/
    cat /proc/cmdline > ${ROOT}${ETCD2B}/logs/cmdline.bootstrap
}

# Boot the new OS immediately
kexec_boot() {
    local kopts=${1}
    local kver=$(cd ${ROOT}/lib/modules; ls | head -1)  # should be just one anyway but let's be sure
    local kernel="vmlinuz-${kver}"
    local initrd="initramfs-${kver}.img"

    # --reset-vga seems to cause more issues than it solves
    kexec --load --initrd=${BOOT}/${initrd} --command-line="${kopts}" ${BOOT}/${kernel}
    umount_fss
    kexec --exec
}

