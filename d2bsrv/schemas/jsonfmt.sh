#!/bin/bash

set -eu

fatal() {
    echo "$1" >&2
    exit 1
}

which python &>/dev/null || fatal "You need "python" in the PATH to run this script"
which jq &>/dev/null || fatal "You need "jq" in the PATH to run this script"

for f in $(ls *.json); do
    echo $f
    python -m json.tool $f >/dev/null
    jq -M . $f > .jsonfmt.swp
    mv .jsonfmt.swp $f
done
