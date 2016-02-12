#!/bin/bash

set -eu

DIR="/var/lib/dock2box/backups"

[ -d "${DIR}" ] || mkdir -p ${DIR}

DATETIME=$(date +'%Y%m%d%H%M%S')
etcdtool -p http://localhost:2379 export / | gzip >${DIR}/${DATETIME}.json.gz
