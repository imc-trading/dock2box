#!/bin/bash

# Info messages
info() {
    local message="$1"

    echo -e "\e[32m+ ${message}\e[0m"
}

# Fatal messages, will print message and reboot
fatal() {
    local message="$1"

    echo -e "\e[31m! ${message}"
    echo
    echo -e "Will wait 60 seconds and then reboot\e[0m"
    sleep 60
    killall -15 tee
    reboot
}

# Get kernel option
get_kopt() {
    local key="$1"

    grep --only-matching "${key}=[^ ]*" "/proc/cmdline" | cut -d= -f2
}

# Get kernel option or return default when it's empty
get_kopt_default() {
    local key="$1"

    grep --only-matching "${key}=[^ ]*" "/proc/cmdline" | cut -d= -f2
}

# Download file
download() {
    local url="$1" target="$2"

    wget --quiet --dns-timeout=2 --connect-timeout=2 --read-timeout=2 --no-check-certificate --output-document "${target}" "${url}"
}

set -eu

BASE="/root"
LOG="${BASE}/bootstrap.log"

# Setup logging
info "Logging to: ${LOG}"
exec &> >(tee -a ${LOG})

# Get options passed by PXE menu
DEBUG=$(get_kopt_default "D2B_DEBUG" "true")
info "Debug: ${DEBUG}"
[[ "${DEBUG}" =~ ^[Tt]rue$ ]] && set -x

SCRIPT=$(get_kopt "D2B_SCRIPT")
[ -n "${SCRIPT}" ] || fatal "Missing kernel option: SCRIPT"
info "Script: ${SCRIPT}"

# Export variables
export DEBUG SCRIPT

# Download and execute script
info "Download script"
download "${SCRIPT}" "${BASE}/install.sh"

info "Execute script"
chmod +x "${BASE}/install.sh"
exec "${BASE}/install.sh"
