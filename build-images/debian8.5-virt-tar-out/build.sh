#!/bin/bash

set -eux

NAME=$1
BUILDER=$2

# Debug
#export LIBGUESTFS_DEBUG=1
#export LIBGUESTFS_TRACE=1

case ${BUILDER} in
    virtualbox-iso)
        virt-tar-out -a packer_output/${NAME}-disk1.vmdk / - | gzip --best > packer_output/${NAME}.tar.gz
    ;;
    vmware-iso)
        virt-tar-out -a packer_output/disk.vmdk / - | gzip --best > packer_output/${NAME}.tar.gz
    ;;
    qemu)
        virt-tar-out -a packer_output/${NAME} / - | gzip --best > packer_output/${NAME}.tar.gz
    ;;
    *)
        echo "Unknown builder: ${BUILDER}"
        exit 1
    ;;
esac
