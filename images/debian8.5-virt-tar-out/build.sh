#!/bin/bash

set -eu

NAME=$1
BUILDER=$2

# Debug
#export LIBGUESTFS_DEBUG=1
#export LIBGUESTFS_TRACE=1

case ${BUILDER} in
    virtualbox-iso)
        virt-tar-out -a packer_output/${NAME}-disk1.vmdk / -
    ;;
    vmware-iso)
        virt-tar-out -a packer_output/disk.vmdk / -
    ;;
    qemu)
        virt-tar-out -a packer_output/${NAME} / -
    ;;
    *)
        echo "Unknown builder: ${BUILDER}"
        exit 1
    ;;
esac
