#!/bin/bash

# Constants
readonly TRUE=0
readonly FALSE=1

# Print boolean
print_bool() {
    local value="$1"

    [ ${value} -eq ${TRUE} ] && echo "true" && return
    echo "false"
}

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

# Get flag from kernel options
get_kopt_flag() {
    local key="$1"

    grep "${key}" /proc/cmdline &>/dev/null && echo ${TRUE} && return
    echo ${FALSE}
}

# Get kernel option
get_kopt() {
    local key="$1"

    grep --only-matching "${key}=[^ ]*" "/proc/cmdline" | cut -d= -f2
}

# Download file
download() {
    local url="$1" target="$2"

    wget --quiet --dns-timeout=2 --connect-timeout=2 --read-timeout=2 --no-check-certificate --output-document "${target}" "${url}"
}

set -eu

# Defaults
BASE="/root"
LOG="${BASE}/bootstrap.log"

# Setup logging
info "Logging to: ${LOG}"
exec &> >(tee -a ${LOG})

# Get options passed by PXE menu
DEBUG=$(get_kopt_flag "D2B_DEBUG")
info "Debug: $(print_bool ${DEBUG})"
[ ${DEBUG} -eq ${TRUE} ] && set -x

SCRIPT=$(get_kopt "D2B_SCRIPT")
[ -n "${SCRIPT}" ] || fatal "Missing kernel option: D2B_SCRIPT"
info "Script: ${SCRIPT}"

DISTRO=$(get_kopt "D2B_DISTRO")
[ -n "${SCRIPT}" ] || fatal "Missing kernel option: D2B_DISTRO"
info "Distro: ${DISTRO}"

# Export variables
export DEBUG SCRIPT DISTRO

# Download scripts and execute install.sh
info "Download scripts"
download "${SCRIPT}" "${BASE}/install.sh"
download "${SCRIPT/install.sh/install_functions.sh}" "${BASE}/install_functions.sh"
download "${SCRIPT/install.sh/install_functions-${DISTRO}.sh}" "${BASE}/install_functions-${DISTRO}.sh"

info "Execute script"
chmod +x "${BASE}/install.sh"
exec "${BASE}/install.sh"
