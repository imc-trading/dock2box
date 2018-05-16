lvm_wipe() {
	local disk=$1

	set +e
	udevadm trigger
	udevadm settle

	pvremove --force --force $disk

	sleep 1s
	pvscan
	set -e
}

lvm_create_pv() {
	local dev=$1

	pvcreate --force --force --yes $dev
}

lvm_create_vg() {
	local dev=$1 vg=$2

	vgcreate $vg $dev
}

lvm_create_lv() {
	local vg=$1 lv=$2 size=$3

	if echo $size | grep -E '^[0-9]+%(VG|PVS|FREE)$'; then
		lvcreate -y -n $lv -l $size $vg
		return
	fi

	lvcreate -y -n $lv -L $size $vg
}

lvm_remove_lv() {
	local vg=$1 lv=$2

	lvremove -y $vg/$lv
}
