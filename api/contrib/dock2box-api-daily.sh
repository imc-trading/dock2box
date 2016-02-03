#!/bin/bash

DIR="/var/lib/dock2box-api/daily"

[ -d "${DAY}" ] || mkdir -p ${DIR}

DAY=$(date +'%d')
etcdtool -p http://localhost:5001 export / >/var/lib/dock2box-api/daily/${DAY}.json
