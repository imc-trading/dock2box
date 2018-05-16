mdraid_wipe() {
	set +e

	for dev in $(
		cd /dev
		ls -d md? 2>/dev/null
	); do
		mdadm --stop ${dev}
		mdadm --remove ${dev}
	done

	udevadm control --reload
	sleep 1

	set -e
}

mdraid_create() {
	local disk=$1 disk2=$2
	local md0="/dev/md0" md1="/dev/md1"
	local devroot=$md0
	local counter=0

	# Set RAID flags
	parted --script $disk set 2 raid on
	parted --script $disk set 3 raid on
	parted --script $disk2 set 2 raid on
	parted --script $disk2 set 3 raid on

	# Set boot flags
	parted --script $disk set 2 boot on
	parted --script $disk2 set 2 boot on

	# Create raid
	yes | mdadm --create ${md0} --level=1 --force --raid-devices=2 --metadata=1.0 ${disk}2 ${disk2}2
	yes | mdadm --create ${md1} --level=1 --force --raid-devices=2 --metadata=1.0 ${disk}3 ${disk2}3
}

mdraid_config() {
	local root=$1

	mdadm --detail --scan >$root/etc/mdadm.conf
}
