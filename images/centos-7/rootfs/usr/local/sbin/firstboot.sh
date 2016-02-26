#!/bin/bash

LOCKFILE='/etc/firstboot.d/firstboot.lock'

[ -e "$LOCKFILE" ] && exit 1
# touch file now to avoid issues with ill-behaved scripts
touch $LOCKFILE

for file in /etc/firstboot.d/*.sh; do
    if [ -x "$file" ]; then
        $file
    fi
done
