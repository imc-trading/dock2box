#!/bin/bash

# Uses the following options passed by PXE menu
#
# D2B_DEBUG             flag
# D2B_SCRIPT_URL        string
# D2B_DISTRO            string

set -eu

# Defaults
BASE="/root"
LOG="${BASE}/bootstrap.log"

# Source functions
source "${BASE}/functions.sh"

# Setup logging
info "Logging to: ${LOG}"
exec &> >(tee -a ${LOG})

# Get options passed by PXE menu
DEBUG=$(get_kopt_flag "D2B_DEBUG")
info "Debug: $(print_bool ${DEBUG})"
[ ${DEBUG} -eq ${TRUE} ] && set -x

SCRIPT_URL=$(get_kopt "D2B_SCRIPT_URL")
[ -n "${SCRIPT_URL}" ] || fatal "Missing kernel option: D2B_SCRIPT_URL"
info "Script: ${SCRIPT_URL}"

DISTRO=$(get_kopt "D2B_DISTRO")
[ -n "${DISTRO}" ] || fatal "Missing kernel option: D2B_DISTRO"
info "Distro: ${DISTRO}"

# Export variables
export DEBUG SCRIPT_URL DISTRO

# Download scripts and execute install.sh
info "Download scripts"
#download "${SCRIPT_URL}/install.sh" "${BASE}/install.sh"
#download "${SCRIPT_URL}/install_functions.sh" "${BASE}/install_functions.sh"
#download "${SCRIPT_URL}/install_functions-${DISTRO}.sh" "${BASE}/install_functions-${DISTRO}.sh"

info "Execute script install.sh"
chmod +x "${BASE}/install.sh"
exec "${BASE}/install.sh"
