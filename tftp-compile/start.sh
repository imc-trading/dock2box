#!/bin/bash

set -eux

compile() {
    server_name="$1"

    sed s/%%server_name%%/${server_name}/g /embed_template.ipxe >/embed.ipxe
    cd /ipxe-${IPXE_VERS}/src && make CFLAGS_config="-DIMAGE_COMBOOT=on" EMBED=/embed.ipxe bin/undionly.kpxe
    cp /ipxe-${IPXE_VERS}/src/bin/undionly.kpxe /var/tftpboot/undionly.kpxe
    chmod 0444 /var/tftpboot/undionly.kpxe
}

SERVER_NAME='${next-server}'

if [[ $# -gt 0 ]]; then
    SERVER_NAME="$1"
    shift
fi

compile "${SERVER_NAME}"

exec /usr/sbin/in.tftpd --foreground -m /var/tftpboot/mapfile --user tftp --secure /var/tftpboot/ "$@"
