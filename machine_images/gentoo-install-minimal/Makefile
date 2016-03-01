# created by Jean-Michel Smith, April 2015
# extracts required files from gentoo minimum-install-cd ISO and builds minimum Dock2box image
# 
# advantages: optimized, reduced size, newer kernel, more robust rescue and bootstrap process

include ../global-d2b-infra.mk

IMGKEY=../shared/id_rsa

NAME=gentoo-install-minimal-generic
ARCH=x86_64
GARCH=amd64
IMAGENAME=stage3-$(GARCH)
IMAGEVERSION=20150618
JSON={"IMAGE_URI": "$(IMAGE_URI)", "Image": "$(IMAGENAME)", "Version": "$(IMAGEVERSION)"}

RELEASE:=$(shell date -u +'%Y%m%d%H%M%S')
TIMESTAMP:=$(shell date +'%Y%m%d-%H%M%S')

TOPDIR:=$(shell pwd)
TMP=$(TOPDIR)/tmp
ISO=$(TMP)/iso
INITRD=$(TMP)/initrd.dir
SQUASHFILE=image.squashfs
SQUASHTMPDIR=$(TMP)/squashfs
CUSTOMDIR=$(TMP)/squashfs-root
# we need ISO and IMAGE to obtain the initramfs
ISO=$(TMP)/iso
IMAGE=install-amd64-minimal-latest.iso

DZDO:=$(shell which dzdo 2> /dev/null)
ifeq ("$(DZDO)", "")
SUDO:=$(shell which sudo)
else
SUDO:=$(DZDO)
endif

KERNEL:=$(shell uname -s)

all: etcd

clean:
	-$(SUDO) umount $(TMP)/$(SYSRCDFILE) $(SQUASHTMPDIR) $(ISO) 
	-$(SUDO) rm -rf $(TMP) logs image.* initrd* initramfs* kernel vmlinuz *.igz *.dat *.md5 *.tar.gz *.tar.bz2 *.zip *.iso $(IMAGE)
	for file in $$(find . -name '*.sed'); do \
		rm -f $${file%%.sed} ; \
	done
	-docker rmi $(NAME) &>/dev/null || true

# obtain the min-install-cd
download: clean
	wget -c $(GENTOO_ISO_URI)
	wget -c $(MEGACLI_URI)

# extract image.squashfs, kernel etc from ISO
extract: download
	-mkdir -p $(TMP) $(ISO) $(SQUASHTMPDIR)
	$(SUDO) mount -o ro,loop $(IMAGE) $(ISO)
	$(SUDO) cp $(ISO)/{image.squashfs,isolinux/gentoo,isolinux/gentoo.igz} $(TMP)
	$(SUDO) umount $(ISO)
	# rename kernel
	mv $(TMP)/gentoo $(TMP)/kernel 

# copy squashfs into tar file so docker can manipulate it
squashfs2tar: extract
	$(SUDO) mount -t squashfs -o loop $(TMP)/$(SQUASHFILE) $(SQUASHTMPDIR)
	$(SUDO) bash -c "cd $(SQUASHTMPDIR) && tar jcpf $(TOPDIR)/image-squashfs.tar.bz2 ."
	$(SUDO) umount $(SQUASHTMPDIR)

# substitue json defined variables in Dockerfile (replacing tf)
# this is not as clean as tf, and requires knowledge of the Dockerfile.sed internals, which tf avoided,
# but it does reduce external dependencies
json_substitutions: squashfs2tar
	cp Dockerfile.sed Dockerfile
	sed -i 's|{{ IMAGE_URI }}|$(IMAGE_URI)|g' Dockerfile
	sed -i 's|{{ Image }}|$(IMAGENAME)|g' Dockerfile
	sed -i 's|{{ Version }}|$(IMAGEVERSION)|g' Dockerfile
	cp rootfs/etc/init.d/pwgen.sed rootfs/etc/init.d/pwgen
	sed -i 's|{{ ENCRYPTED_ROOT_PASSWORD }}|$(ENCRYPTED_ROOT_PASSWORD)|g' rootfs/etc/init.d/pwgen

# build the docker image
dockerimage: json_substitutions
	mkdir logs
	docker build --no-cache=true --rm=true -t $(NAME):latest . | tee logs/dockerbuild.log

# get docker image into squashfs staging directory
docker2squashfs: dockerimage
	mkdir -p $(CUSTOMDIR) $(INITRD)
	$(SUDO) docker run -i -t --entrypoint /bin/true $(NAME):latest
	$(SUDO) docker export $$(docker ps -a -l | awk '/$(NAME)/ {print $$1}') | $(SUDO) tar xpf - -C $(CUSTOMDIR)

# extract kernel
kernel: docker2squashfs
	$(SUDO) cp $(CUSTOMDIR)/boot/kernel* $(TMP)/kernel

# extract initramfs
initramfs: kernel
	# stage directory and download image
	-mkdir -p $(ISO)
	# this worked, oddly enough:
	#wget -c $(IMAGE_URI)/$(IMAGE)
	# but we should probably be using:
	wget -c $(GENTOO_ISO_URI)
	# mount image and extract initramfs
	$(SUDO) mount -o ro,loop $(IMAGE) $(ISO)
	$(SUDO) bash -c "xz -dc $(ISO)/isolinux/gentoo.igz | ( cd $(INITRD) && cpio -idv )"
	# unmount as soon as possible to keep things clean
	$(SUDO) umount $(ISO)

# patch initramfs and put kernel modules in place
build: initramfs
	# patch init
	$(SUDO) patch -d $(INITRD) -p0 < init.patch
	# put kernel modules in place
	$(SUDO) cp -a $(CUSTOMDIR)/lib/modules $(INITRD)/lib/modules
	# adjust location of /usr/portage to point to a non-read-only location so we can emerge on the fly if we need to:
	$(SUDO) rm -rf $(CUSTOMDIR)/usr/portage
	$(SUDO) ln -s /portage $(CUSTOMDIR)/usr/portage

# now we're ready to compress squashfs
squashfs: build
    # patch fstab
	$(SUDO) patch -d $(CUSTOMDIR)/etc -p0 < fstab.patch
	-rm -f $(TMP)/$(SQUASHFILE)
	$(SUDO) bash -c "cd $(TMP) && mksquashfs $(CUSTOMDIR) $(SQUASHFILE) -noappend -always-use-fragments"

# combine squashfs.image with initrd
combine: squashfs
	# Now add our docker-created "squashfs" to it,
	$(SUDO) mkdir -p "$(INITRD)/mnt/cdrom"
	$(SUDO) cp "$(TMP)/$(SQUASHFILE)" "$(INITRD)/mnt/cdrom"

# and build the initrd containing the revised squash.image
initrd: combine
	$(SUDO) bash -c "(cd $(INITRD) && find . -print | cpio -o -H newc | gzip -9 -c - ) > $(TMP)/initrd"
	mv -f $(TMP)/kernel $(TMP)/initrd .

# publish this baby. 
publish: initrd 
	chmod 600 $(IMGKEY)
	ssh -q -o 'UserKnownHostsFile /dev/null' -o 'StrictHostKeyChecking no' -i $(IMGKEY) $(IMGUSR)@$(IMGSRV) mkdir -p $(IMGPTH)/$(NAME)-$(ARCH)-$(RELEASE)
	rsync --exclude=*raw -av -e 'ssh -q -i $(IMGKEY) -o "UserKnownHostsFile /dev/null" -o "StrictHostKeyChecking no"' --delete --max-delete=10 \
		--chmod=Dug=rw,Do=rx,Fug=rw,Fo=r bootstrap.sh functions.sh kernel initrd $(IMGUSR)@$(IMGSRV):$(IMGPTH)/$(NAME)-$(ARCH)-$(RELEASE)
	ssh -q -i $(IMGKEY) -o 'UserKnownHostsFile /dev/null' -o 'StrictHostKeyChecking no' $(IMGUSR)@$(IMGSRV) "cd $(IMGPTH); \
		rm -f $(NAME)-$(ARCH)-latest; ln -sv $(NAME)-$(ARCH)-$(RELEASE) $(NAME)-$(ARCH)-latest"

# use this if you want to also store the image in the docker registry (will rebuild docker image accordingly)
push: squashfs2tar
	docker login -u $(USER) -p '$(PASS)' -e $(MAIL) $(REG)
	docker build --no-cache --rm=true -t $(REG)/$(NAME):$(REL) .
	docker tag -f $(REG)/$(NAME):$(REL) $(REG)/$(NAME):latest
	#docker login -u $(USER) -p '$(PASS)' -e $(MAIL) $(REG)
	docker push $(REG)/$(NAME):$(REL)
	docker push $(REG)/$(NAME):latest

etcd: publish
	etcd-cli -y set image --image $(NAME)-x86_64 --file_registry $(IMGPTHURL) --docker_registry $(REG) --kopts "biosdevname=0 net.ifnames=0 looptype=squashfs cdroot=1 init=/linuxrc real_root=/ setkmap=us root=/dev/ram0 loop=/image.squashfs" --bootstrap_image $(NAME)-x86_64-latest --timestamp=$(TIMESTAMP)

