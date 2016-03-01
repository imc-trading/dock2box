# Base image
FROM scratch
ADD image-squashfs.tar.bz2 /

# Maintainer/author
MAINTAINER Jean-Michel Smith <jean@jean-michel.eu>

# BOOTSTRAP-BASE
RUN cp -a /usr/portage /usr/portage.clean && \ 
    cp /etc/fstab /etc/fstab.orig && \ 
    wget -q -c {{ IMAGE_URI }}/{{ Image }}-{{ Version }}.tar.bz2 && \ 
    tar jxpf {{ Image }}-{{ Version }}.tar.bz2 --exclude='./sys' --exclude='./dev' --exclude='./etc/hosts' && \ 
    rm {{ Image }}-{{ Version }}.tar.bz2 && \ 
    emerge-webrsync ; eselect news read 

RUN mkdir -p /etc/portage/package.mask /etc/portage/package.unmask /etc/portage/package.use \
        /etc/portage/package.accept_keywords /etc/portage/package/env /etc/portage/package/sets \
        /etc/portage/repos.conf 

COPY rootfs/etc/portage/package.mask/ /etc/portage/package
COPY rootfs/etc/portage/make.conf /etc/portage/make.conf
COPY rootfs/etc/locale.gen /etc/locale.gen
COPY rootfs/etc/portage/package.accept_keywords/ /etc/portage/package.accept_keywords/
COPY rootfs/etc/portage/package.use/ /etc/portage/package.use/

# Seems early, but this fixes python, openssl, and a host of other issues:
RUN emerge -DuvN --with-bdeps=y @world && \
    emerge @preserved-rebuild && \
    emerge app-misc/c_rehash && \
    emerge app-admin/python-updater && \
    python-updater && \
    emerge dev-vcs/git && \
    yes | etc-update --automode -3 /etc 
COPY rootfs/etc/portage/repos.conf/* /etc/portage/repos.conf/
RUN rm -rf /usr/portage ; mkdir /portage ; ln -s /usr/portage /portage && \
    # Break build on non-Gentoo because Gentoo is the only distro with a reliable console ; \
    #git config --global http.sslVerify false ; \
    emerge --sync 

# You should add any in-house or custom certs here:
# place your cert files in rootfs/usr/share/ca-certificates/<yourfirmname>/ and uncomment the following
# (e.g. change "myfirm/" to something with meaning for you, and filenames to whatever your cert files are called)
#RUN mkdir -p /usr/share/ca-certificates/myfirm && \
#    echo "myfirm/ca.crt" >> /etc/ca-certificates.conf && \
#    echo "myfirm/catrading01.pem" >> /etc/ca-certificates.conf && \
#    echo "myfirm/catrading02.pem" >> /etc/ca-certificates.conf && \
#    echo "myfirm/catrading03.pem" >> /etc/ca-certificates.conf
#COPY rootfs/usr/share/ca-certificates/myfirm/ /usr/share/ca-certificates/myfirm/
#RUN update-ca-certificates

# install layman
RUN emerge app-portage/layman ; \
    echo "add your layman overlays here.  E.g. 'layman -a blah' or, for unoffical overlays, 'layman -o http://myhttphost/full/path/to/my-overlay/overlay.xml -a my-overlay'"
# add this if you do use a layman-managed overlay:
#RUN echo "source /var/lib/layman/make.conf" >> /etc/portage/make.conf 

# and emerge basic stuff
# pam is needed for util-linux to compile
RUN emerge app-editors/vim && \ 
    eselect vi set vim && \ 
    eselect editor set /usr/bin/vi && \ 
    emerge sys-process/htop && \
    emerge net-misc/ipcalc && \
    emerge app-portage/cpuinfo2cpuflags && \ 
    emerge udev && \
    emerge memtest86+ gentoolkit dmraid lvm2 livecd-tools sys-fs/mdadm && \ 
    emerge scripts mingetty && \
    emerge sys-apps/lsb-release && \
    emerge sys-libs/pam && \
    emerge sys-apps/util-linux

# make sure JQ is installed
RUN emerge sys-devel/bison && \
    emerge app-misc/jq 

# MegaCLI
COPY 8.07.10_MegaCLI_Linux.zip /distfiles/8.07.10_MegaCLI_Linux.zip
# MegaCLI, aspell-en (for EMEA-style auto-hostname generation), etc.
RUN emerge app-arch/rpm2targz dev-db/etcd && \ 
    emerge sys-block/megacli net-misc/lldpd sys-apps/biosdevname net-analyzer/fping sys-apps/kexec-tools sys-apps/ipmitool && \ 
    emerge app-text/aspell app-dicts/aspell-en 

# REBUILD WORLD TO ENSURE EVERYTHING IS IN A GOOD STATE
RUN emerge -Deuv --with-bdeps=y @world && \
    yes | etc-update --automode -3 /etc 

# set up startup scripts
COPY rootfs/etc/init.d/pwgen /etc/init.d/pwgen

# set your server sshd defaults as appropriate for your environment.  This will allow remote access to a troublesome host booted with the bootstrap image
# which is very useful when debugging:
RUN rc-update add sshd default && \
    sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/g' /etc/ssh/sshd_config && \
    sed -i 's/#PermitRootLogin yes/PermitRootLogin yes/g' /etc/ssh/sshd_config && \
    sed -i 's/# PermitRootLogin yes/PermitRootLogin yes/g' /etc/ssh/sshd_config && \
    sed -i 's/localhost/dock2box-bootstrap/g' /etc/conf.d/hostname

# BOOTSTRAP-BASE-KERNEL

RUN emerge sys-kernel/genkernel sys-kernel/gentoo-sources 

# build new kernel
COPY rootfs/usr/src/config-4.4 /usr/src/linux/.config
RUN rm -rf /lib/modules/* && \ 
    bash -c "cd /usr/src/linux && make oldconfig && genkernel --oldconfig --no-splash --no-clean all && make clean"

# BOOTSTRAP-BASE-DOCKER

# This uses LVM as the docker backend for bootstrap purposes.  This is handy because we can explicity set the image size large enough to handle
# whatever we need:
RUN emerge app-text/docbook-xml-dtd app-text/docbook-xsl-stylesheets app-text/docbook-xsl-ns-stylesheets && \ 
    emerge sys-fs/btrfs-progs && \ 
    emerge app-emulation/docker && \
    sed -i 's/DOCKER_OPTS=""/DOCKER_OPTS="--storage-opt dm.basesize=40G --storage-driver=devicemapper"/g' /etc/conf.d/docker

# BOOTSTRAP FINAL

# add dictionaries used by bootstrap.sh to image (e.g. for hands-free host auto-naming)
COPY rootfs/dictionaries /

# add the auto-start bootstrap glue and rename the image hostname to 'dock2box-bootstrap'
COPY rootfs/etc/init.d/dock2box-bootstrap /etc/init.d/dock2box-bootstrap
COPY rootfs/etc/motd /etc/motd
RUN sed -i 's/use_lvmetad = 0/use_lvmetad = 1/g' /etc/lvm/lvm.conf && \
    rc-update add dock2box-bootstrap default 

# CLEANUP

# restore fstab, remove kernel sources, clean up portage and revert to "portage.clean"
RUN mv /etc/fstab.orig /etc/fstab && \
    find /etc -name '._*' -exec rm {} \; && \ 
    find /usr/share/man -type f -exec rm {} \; && \ 
    rm -f /usr/share/info/* && \ 
    rm -rf /var/tmp/portage/* /tmp/* /var/cache/edb/dep && \
    emerge -C gentoo-sources && rm -rf /usr/src/* && \ 
    emerge app-admin/localepurge && \
    echo "Removing portage tree.  This is done only after last emerge command is run, to minimize image size." && \
    rm -rf /distfiles/* && \
    echo "MANDELETE" > /etc/locale.nopurge && \ 
    echo "SHOWFREEDSPACE" >> /etc/locale.nopurge && \ 
    echo "VERBOSE" >> /etc/locale.nopurge && \ 
    echo "en" >> /etc/locale.nopurge && \ 
    echo "en_US" >> /etc/locale.nopurge && \ 
    echo "en_US.UTF-8" >> /etc/locale.nopurge && \ 
    echo "en_US.iso88591" >> /etc/locale.nopurge && \ 
    localepurge && \ 
    rm -rf /portage /usr/portage.clean && \
    find / -path /proc -prune -path /dev -prune -type l ! -xtype f ! -xtype d -ok rm -f {} \; && \ 
    find / -path /proc -prune -path /dev -prune -type f -xdev -name ".keep" -print -exec rm {} \;

