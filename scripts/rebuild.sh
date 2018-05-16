#!/bin/bash

set -eux

export PATH=/usr/bin:/bin:/usr/sbin:/sbin

fatal() {
	echo "FATAL: $1" >&2
	exit 1
}

# Set defaults
[ -n "${KEXEC:-}" ] || KEXEC="false"
[ -n "${DEBUG:-}" ] || DEBUG="false"
[ -n "${KARGS:-}" ] || KARGS="video=800x600-32 quiet splash"
[ -n "${KARGS_DEBUG:-}" ] || KARGS_DEBUG="video=800x600-32 debug"
[ -n "${KERNEL:-}" ] || KERNEL="http://dock2box/boot/kernel"
[ -n "${INITRD:-}" ] || INITRD="http://dock2box/boot/initrd"

# Download kernel and initrd
curl -s $KERNEL -o /boot/dock2box-kernel
curl -s $INITRD -o /boot/dock2box-initrd

# Template Grub entry
cat <<EOF >/etc/grub.d/50_dock2box
#!/bin/sh
cat <<MENU
menuentry 'Dock2Box' {
        linux16 /dock2box-kernel root=\${GRUB_DEVICE} ${KARGS}
        initrd16 /dock2box-initrd
}

menuentry 'Dock2Box Debug' {
        linux16 /dock2box-kernel root=\${GRUB_DEVICE} ${KARGS_DEBUG}
        initrd16 /dock2box-initrd
}
MENU
EOF
chmod 755 /etc/grub.d/50_dock2box

# Create Grub entry
grub2-mkconfig -o /boot/grub2/grub.cfg

# Set Grub default entry
grub2-set-default Dock2Box

# Reboot using Kexec
if [ "${KEXEC:-}" == "true" ]; then
	if [ "${DEBUG:-}" == "true" ]; then
		kexec -l /boot/dock2box-kernel --initrd /boot/dock2box-initrd --command-line "root=/dev/mapper/vg_root-lv_root ${KARGS}"
	else
		kexec -l /boot/dock2box-kernel --initrd /boot/dock2box-initrd --command-line "root=/dev/mapper/vg_root-lv_root ${KARGS_DEBUG}"
	fi
fi
