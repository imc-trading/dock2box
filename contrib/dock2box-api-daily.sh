#!/bin/bash

set -eu

DIR="/var/lib/dock2box-api/daily"

log() {
    logger -s -t $(basename $0) $1
}

[ -d "${DIR}" ] || mkdir -p ${DIR}

log "Started"

DAY=$(date +'%d')
etcdtool -p http://localhost:5001 export / >/var/lib/dock2box-api/daily/${DAY}.json

log "Done"
