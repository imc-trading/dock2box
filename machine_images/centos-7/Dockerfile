# Base image
FROM centos:centos7

# Maintainer/author
MAINTAINER Michael Persson, Adrian Vasile, Jean-Michel Smith

# Install SystemD
# centos 7.0, 7.1:
RUN yum -y swap -- remove fakesystemd -- install systemd systemd-libs
# centos >=7.2?
# RUN yum -y swap -- remove systemd-container systemd-container-libs -- install systemd systemd-libs

# Dracut config, needs to be before the kernel install
COPY rootfs/etc/dracut.conf /etc/dracut.conf

# Cleanup weird ami file that showed up in the latest upstream image, breaks our custom driver list
RUN rm -f /etc/dracut.conf.d/ami.conf

# Install packages
RUN yum install -y epel-release; \
    yum install -y @base net-tools bind-utils mailx bzip2 \
        unzip wget git wget curl dmidecode sharutils \
        dracut grub2 grub2-tools biosdevname e2fsprogs \
        kernel kernel-devel kernel-tools openssh-server \
        grub2-efi shim rsyslog parted gdisk \
        net-tools bind-utils mailx bzip2 unzip wget \
        git wget curl dmidecode sharutils sudo \
        redhat-lsb python-devel python-pip gcc jq ipmitool \
        openssh-server openssh-clients

# Configure timezone
RUN ln -sf /usr/share/zoneinfo/Etc/UTC /etc/localtime
COPY rootfs/etc/sysconfig/clock /etc/sysconfig/clock

# Configure ntpdate
#RUN yum install -y ntpdate 
#RUN /usr/bin/systemctl enable ntpdate
#COPY rootfs/etc/sysconfig/ntpdate /etc/sysconfig/ntpdate

# Remove chrony since it prevents ntp from starting at boot
#RUN yum remove -y chrony

# Configure ntp
#RUN yum install -y ntp 
#RUN /usr/bin/systemctl enable ntpd
#COPY rootfs/etc/ntp.conf /etc/ntp.conf

# Disable selinux and install deps
RUN yum install -y libselinux-utils libselinux selinux-policy-targeted \
    selinux-policy libselinux-python
COPY rootfs/etc/selinux/config /etc/selinux/config

# Set root password to 'dock2box'
RUN usermod -p '$6$WyzcP9uG$64NvhY1P.W1d08CRxFjJ5QUDyZcrzyjxVPV82bAxBgf9eE2sdSZs.v47HGdWPvLxhNxAkrJten86R2Qb1vdhe/' root

# Add firstboot
COPY rootfs/usr/local/sbin/firstboot.sh /usr/local/sbin/firstboot.sh
RUN chmod +x /usr/local/sbin/firstboot.sh ;\
    mkdir -p /etc/firstboot.d
COPY rootfs/etc/systemd/system/firstboot.service /etc/systemd/system/firstboot.service
RUN systemctl enable firstboot
COPY rootfs/etc/firstboot.d/* /etc/firstboot.d/
RUN chmod 755 /etc/firstboot.d/*.sh

# Install Peekaboo
# Install from GitHub binary RPM file
#RUN yum install -y peekaboo ;\
#    systemctl enable peekaboo

# Add local IP to Console
# Could be improved, assumes eth0 is the interface
COPY rootfs/sbin/ifup-local /sbin/ifup-local
RUN chmod +x /sbin/ifup-local

# Fix issue with unresponsive Console:
# tty_ldisc_hangup_halt: wating (login) for tty1 took too long, but we keep waiting…
RUN echo "xvc0" >> /etc/securetty

# Allow users in group wheel to use sudo
COPY rootfs/etc/sudoers.d/wheel /etc/sudoers.d/wheel

# Cleanup yum
RUN yum clean all

# Cleanup log files
RUN find /var/log -type f | while read f; do echo -ne '' > $f; done;

# Cleanup under tmp directory
#RUN rm -rf /tmp/* /etc/sysconfig/network-scripts/ifcfg-default
