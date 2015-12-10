#/bin/bash

set -eu

cpt() {
    printf "\n########## $1 ##########\n\n"
}

fatal() {
    local msg="$1"

    echo "$msg" >&2
    exit 1
}

# Check for pre. requisites
which mongo &>/dev/null || fatal "Missing pre. requisite: mongo"

#
# Drop database
#
cpt "Drop Database"
printf "use d2b\ndb.dropDatabase()\n" | mongo
