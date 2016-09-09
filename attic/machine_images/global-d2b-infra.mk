# Author: Jean-Michel Smith, Feb 25, 2016
#
# This file contains definitions specific to each organization's Dock2Box infrastructure, defining where you run your http daemon/fileserver, your
# Docker registry, and so on.

# *** HTTP / IMAGE SERVER, USER, and LOCATION DEFINITIONS ***

# Set this to the encrypted root password you want for your bootstrap image. This is the string exactly as it appears in /etc/shadow
# generic default (CHANGE THIS) is 'Welcome@d2b'
ENCRYPTED_ROOT_PASSWORD := '$6$9A.1yxCC$ZFkiRXzf2pdvstshog17xBjcCVAQwVtqtwsn6PTt63yBNiKS9ywHMKr3XUZ2e7SXbm1fCwmXY9AQCVvDzOkAR.'

# Set this to the location of your Dock2Box ssh keys (this is a local file.  Often this is easiest if kept in your Dock2Box git repo):
IMGKEY := ../shared/id_rsa

# Set this to the hostname of your HTTP server:
IMGSRV := d2bhttpservername.your.org.com

# Set this to the URI path where your source vanilla images reside:
IMAGE_URI := http://$(IMGSRV)/archives/linux/gentoo/releases/amd64/autobuilds/current-install-amd64-minimal/

# Set this to URI where the gentoo minimal install ISO resides:
GENTOO_ISO_URI := http://$(IMGSRV)/archives/linux/gentoo/releases/amd64/autobuilds/20151231/install-amd64-minimal-latest.iso

# Set this to the URI where the megacli tarball resides:
MEGACLI_URI := http://$(IMGSRV)/archives/linux/sysrescuecd/8.07.10_MegaCLI_Linux.zip

# Set this to the username on your HTTP server Dock2Box should use to download images.
# NOTE: The Dock2Box public ssh key (e.g. ../shared/id_rsa.pub) needs to be in this user's ~/.ssh/authorized_keys file for passwordless scp/rsync to work:
IMGUSR := d2buser

# Set this to the OS Filesystem Path on your HTTP Server where Dock2Box generated images are uploaded to and served from.
# This is not the same path that houses your vanilla source images.  Remember:
# "Vanilla Source Image" (e.g. created with Packer) -> Dock2Box "make image(s)" -> Docker -> D2B Image -> scp/rsync to $(IMGPTH)
IMGPTH := /var/yum/osimages

# This is the full "URI" path where the Dock2Box-generated OS images reside.  Leave as is to use the default location based
# on the values you defined above:
IMGPTHURL=http://$(IMGSRV)/osimages

# *** DOCKER REGISTRY DEFINITIONS ***

# Set this to the account name of your Docker Registry User:
SER=registry_svc

# Set this to the password required for access to your Docker Registry:
PASS=Welcome@d2b

# Set this to the email address of the Docker Registry user:
MAIL=noreply@your.host.com

# Set this to the hostname:port where your Docker Registry resides:
REG=docker-registry.your.org.com:8080

define jsonmerge 
$(shell echo '$1 $2' | jq -s '.[0] * .[1]')
endef

