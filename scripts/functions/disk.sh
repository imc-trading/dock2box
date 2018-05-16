disk_get_first() {
	[ -b "/dev/vda" ] && echo "/dev/vda" && return
	[ -b "/dev/sda" ] && echo "/dev/sda" && return
	[ -b "/dev/nvme0n1" ] && echo "/dev/nvme0n1" && return
	echo ""
}

disk_get_second() {
	dev=""
	[ -b "/dev/sdb" ] && dev="/dev/sdb"
	[ -b "/dev/vdb" ] && dev="/dev/vdb"
	[ -b "/dev/nvme1n1" ] && dev="/dev/nvme1n1"

	if [ -z "${dev}" ]; then
		echo ""
		return
	fi

	if echo q | fdisk $dev &>/dev/null; then
		echo $dev
		return
	fi

	echo ""
}

disk_part_name() {
	local disk=$1

	[ $disk == "/dev/nvme0n1" ] && echo "/dev/nvme0n1p" && return
	[ $disk == "/dev/nvme1n1" ] && echo "/dev/nvme1n1p" && return
	echo $disk
}

disk_size_mb() {
	local disk=$1

	parted --script $disk unit MB print | grep "^Disk $disk" | awk '{print $NF}'
}

disk_mbr_wipe() {
	local disk=$1

	dd if=/dev/zero of=$disk bs=512 count=1
}

disk_part_wipe() {
	local disk=$1

	wipefs --all $disk
}

disk_part_gpt() {
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

	# IMPORTANT: Ensure bios that can't handle GPT directly can still boot
	parted --script --align optimal $disk disk_toggle pmbr_boot
}

disk_part_mbr() {
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

disk_part_boot() {
	local dev=$1

	mkfs.vfat -F32 -n BOOTFS $dev
}

disk_sync() {
	sync
}
