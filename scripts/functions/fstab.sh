fstab_create() {
	local root=$1

	[ -d $root/etc ] || mkdir -p $root/etc

	touch $root/etc/fstab
	fstab_add $root LABEL=rootfs / ext4 defaults 1 1
	fstab_add $root LABEL=BOOTFS /boot vfat defaults 1 2
	fstab_add $root LABEL=swapfs swap swap defaults 0 0
	fstab_add $root tmpfs /tmp tmpfs noatime 0 0
}

fstab_add() {
	local root=$1 name=$2 dir=$3 fs=$4 opts=$5 freq=$6 passno=$7

	printf "%-15s %-15s %-5s %-15s %-1s %-1s\n" $name $dir $fs $opts $freq $passno >>$root/etc/fstab
}
