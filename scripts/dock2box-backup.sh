#!/bin/bash

set -eu

DIR="/var/lib/dock2box/backups"

[ -d "${DIR}" ] || mkdir -p ${DIR}

DATETIME=$(date +'%Y%m%d%H%M%S')
etcdtool -p http://localhost:2379 export / >${DIR}/${DATETIME}.json

if test -e ${DIR}/last.json && diff ${DIR}/last.json ${DIR}/${DATETIME}.json; then
    rm -f ${DIR}/${DATETIME}.json
    echo "No new content, won't keep backup"
    exit 0
fi

rm -f ${DIR}/last.json
ln -s ${DIR}/${DATETIME}.json ${DIR}/last.json
