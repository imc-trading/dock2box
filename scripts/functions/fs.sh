fs_create() {
	local fs=$1 dev=$2 name=$3

	case $fs in
		ext4)
			mkfs.ext4 -F -i 4096 -q -L $name $dev
			;;
		xfs)
			mkfs.xfs -f -L $name $dev
			;;
		swap)
			mkswap -L $name $dev
			;;
	esac
}

fs_mount_get() {
	local dev=$1

	awk "\$2==\"$dev\" {print \$0}" /proc/mounts
}

fs_mounts() {
	local path=$1

	awk "\$2~\"^$dir\" {print \$2}" /proc/mounts
}

fs_mount() {
	local fs=$1 dev=$2 dir=$3

	[ -d $dir ] || mkdir -p $dir

	case $fs in
		ext4 | vfat)
			mount -t $fs $dev $dir
			;;
		bind)
	        [ -d $dev ] || mkdir -p $dev
			mount -o bind $dev $dir
			;;
	esac
}

fs_umount() {
	local dir=$1

	[ -n "$(fs_mount_get $dir)" ] && umount $dir
}

fs_path_depth() {
	local paths="$1"

	for m in $paths; do
		c=$(echo $m | tr -dc / | wc -c)
		echo $c $m
	done
}

fs_path_depth_sort() {
	local paths="$1"

	fs_path_depth "$paths" | sort -r -n | awk '{print $2}'
}

fs_umount_recurse() {
	local dir=$1

	for m in $(fs_path_depth_sort "$(fs_mounts $dir)"); do
		fs_umount $m
	done
}
