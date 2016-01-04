#!/bin/bash

set -eu

fatal() {
    local msg="$1"

    echo "$msg" >&2
    exit 1
}

# Check for pre. requisites
which python &>/dev/null || fatal "Missing pre. requisite: python"
which jq &>/dev/null || fatal "Missing pre. requisite: jq"

# Valiadate and format JSON files
for f in `ls *.json`; do
    echo $f
    python -m json.tool $f >/dev/null
    jq -M . $f > .jsonfmt.swp
    mv .jsonfmt.swp $f
done
