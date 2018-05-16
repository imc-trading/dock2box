install_grub() {
	local root=$1 disks="$2" kopts="$3" mdraid="$4"

	mkdir -p $root/etc/default
	cat <<EOF >$root/etc/default/grub
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR="CentOS Linux"
GRUB_DEFAULT=saved
GRUB_DISABLE_SUBMENU=true
GRUB_TERMINAL_OUTPUT="console"
GRUB_CMDLINE_LINUX="rd.auto rd.auto=1 domdadm dolvm vconsole.keymap=us verbose nomodeset crashkernel=auto selinux=0 pcie_aspm=off biosdevname=0 net.ifname=0${kopts}"
GRUB_DISABLE_RECOVERY="true"
EOF

	local i=0
	true >$root/boot/grub2/device.map
	for disk in $disks; do
		echo "(hd${i})    ${disk}" >>$root/boot/grub2/device.map
		i=$((i + 1))
	done

	if [ "${mdraid:-}" == "true" ]; then
		chroot $root mdadm --wait /dev/md0 || true
		sed -i 's/#add_dracutmodules+=""/add_dracutmodules+="mdraid lvm qemu"/g' $root/etc/dracut.conf
	else
		sed -i 's/#add_dracutmodules+=""/add_dracutmodules+="lvm qemu"/g' $root/etc/dracut.conf
	fi

	chroot $root grub2-mkconfig -o /boot/grub2/grub.cfg

	for disk in $disks; do
		chroot $root grub2-install $disk
	done

	kver=$(ls ${ROOT}/lib/modules/ | tail -1)
	sed -i 's/#hostonly="yes"/hostonly="no"/g' $root/etc/dracut.conf
	sed -i 's/#filesystems+=""/filesystems+="ext4 xfs btrfs"/g' $root/etc/dracut.conf
	sed -i 's/#add_drivers+=""/add_drivers+="vmw_balloon vmw_vmci vmw_pvscsi 3w-9xxx 3w-sas aacraid acard-ahci ahci ahci_platform aic79xx arcmsr ata_generic ata_piix be2iscsi bfa bnx2fc bnx2i ch cxgb3i cxgb4i fcoe fnic hpsa hptiop hv_storvsc initio isci iscsi_boot_sysfs iscsi_tcp libahci libata libcxgbi libfc libfcoe libiscsi libiscsi_tcp libosd libsas libsrp lpfc megaraid_sas mpt2sas mpt3sas mvsas mvumi osd osst pata_acpi pata_ali pata_amd pata_arasan_cf pata_artop pata_atiixp pata_atp867x pata_cmd64x pata_cs5536 pata_hpt366 pata_hpt37x pata_hpt3x2n pata_hpt3x3 pata_it8213 pata_it821x pata_jmicron pata_marvell pata_netcell pata_ninja32 pata_oldpiix pata_pdc2027x pata_pdc202xx_old pata_piccolo pata_rdc pata_sch pata_serverworks pata_sil680 pata_sis pata_via pdc_adma pm80xx pmcraid qla2xxx qla4xxx raid_class sata_mv sata_nv sata_promise sata_qstor sata_sil24 sata_sil sata_sis sata_svw sata_sx4 sata_uli sata_via sata_vsc scsi_debug scsi_tgt scsi_transport_fc scsi_transport_iscsi scsi_transport_sas scsi_transport_spi scsi_transport_srp sd_mod ses sg sr_mod stex st ufshcd ufshcd-pci virtio_scsi vmw_pvscsi raid0 raid1 raid10 raid456 xenblk xen_blkfront virtio_net virtio_pci virtio_blk virtio_balloon e1000 mptbase mptscsih mptsas mptspi i2c-piix4"/g' $root/etc/dracut.conf

	chroot $root dracut -v --force --mdadmconf /boot/initramfs-${kver}.img $kver
}
