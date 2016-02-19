#!/bin/bash

set -eu
#set -x

API_URL="http://dock2box:8080/api/v1"

# Get kernel option
get_kopt() {
    local key="$1"
    grep --only-matching "${key}=[^ ]*" "/proc/cmdline" | cut -d= -f2
}

# Get resource from API
get_resource() {
    local resource="$1"

    response=$(curl -s "${API_URL}${resource}?envelope=true")
    code=$(echo "$response" | jq ".code")
    [ "${code}" != "200" ] && return
    echo "${response}" | jq ".data"
}

# Get key from resource
get_key() {
    local data="$1" key="$2"
    value=$(echo "${data}" | jq -r "${key}")
    [ "${value}" == "null" ] && value=""
    echo "${value}"
}

# Get key from resource or return default when it's empty
get_key_default() {
    local data="$1" key="$2" default="$3"
    value=$(echo "${data}" | jq -r "${key}")
    [ "${value}" == "null" ] && value="${default}"
    echo "${value}"
}

WGET_OPTS="--quiet --dns-timeout=2 --connect-timeout=2 --read-timeout=2 --no-check-certificate"

HWADDR=$(get_kopt "HWADDR")
#DEFAULT_ETH=$(grep -i --with-filename ${BOOT_HWADDR} /sys/class/net/*/address | cut -d/ -f5)
DEFAULT_ETH=eth0
#BOOT_IP=$(get_kopt "IP")
#BOOT_NETMASK=$(get_kopt "NETMASK")
#BOOT_GATEWAY=$(get_kopt "GATEWAY")
#BOOT_SUBNET=$(get_kopt "SUBNET")
IMAGE=$(get_kopt "IMAGE")
IMAGE_VERSION=$(get_kopt "IMAGE_VERSION")
#BOOT_IMAGE_TYPE=$(get_kopt "IMAGE_TYPE")
#BOOT_HOSTNAME=$(get_kopt "HOSTNAME")

# Get site configuration
SITE=$(get_kopt "SITE")
SITE_DATA=$(get_resource "/sites/${SITE}")
DOMAIN=$(get_key "${SITE_DATA}" ".domain")
DNS_SERVERS=$(get_key "${SITE_DATA}" ".dns_servers")

# Get host configuration
FQDN=$(get_kopt "FQDN")
HOST_DATA=$(get_resource "/hosts/${FQDN}")

DEBUG=$(get_key "${HOST_DATA}" ".debug")
if [[ "${DEBUG}" =~ "[Tt]rue" ]]; then
    set -x
fi

# If dhcp is not set default to false
DHCP=$(get_key_default "${HOST_DATA}" ".interface.${BOOT_ETH}.dhcp" "false")

# Get network configuration
if [[ "${DHCP}" =~ ^[Ff]alse$ ]]; then
    IPV4=$(get_key "${HOST_DATA}" ".interface.${BOOT_ETH}.ip")
    NETMASK=$(get_key "${HOST_DATA}" ".interface.${BOOT_ETH}.netmask")
    GATEWAY=$(get_key "${HOST_DATA}" ".interface.${BOOT_ETH}.gateway")

    # If there is no network configuration default to using dhcp
    [ -z "${IPV4}" ] && DHCP="true"
fi

# Get image configuration
IMAGE_DATA=$(get_resource "/hosts/${IMAGE}")
DISTRO=$(get_key "${IMAGE_DATA}" ".distro")

# Not done refactoring
exit 0

# Get default configuartion
#SCRIPT=$(get_key "${DEFAULT}" ".install_script")

# Set initial root password
#/usr/sbin/usermod -p '...' root

# Print cmdline
#grep BOOT --with-filename /proc/cmdline

# Print network information
#ifconfig "${BOOT_ETH}"

# Maximize root filesystem size to appease docker IFF we have >8GB of RAM
MEM_SIZE=$(free -m | grep Mem | awk '{print $2}')

# Strip off last the chars (just show GB)
MEM_GB=$(echo "${MEM_SIZE}" | sed "s/...$//")

if [ ${MEM_GB} -gt 6 ]; then
    TGT_SIZE=$((${MEM_GB} - 2))
    echo "resize / to appropriate size"
    echo "mount -o remount size=${TGT_SIZE}G /"
    mount -o "remount,size=${TGT_SIZE}G" /
fi

# Set hostname
hostname "${FQDN}"

info "Download install script"
download "${SCRIPT}" "/root/install.sh"

echo "Executing the install script"
chmod +x /root/install.sh
export IPV4 NETMASK SUBNET GATEWAY FQDN HOSTNAME
#export IP NETMASK SUBNET GATEWAY HOSTNAME HOSTNAME_SHORT NAMING_SCHEME DOMAIN BOOT_MAC BOOT_ETH DEBUG IMAGE SCRIPT BOOT_IMAGE_NAME BOOT_IMAGE_VERSION BOOT_IMAGE_TYPE DNS

/root/install.sh

CURRENT_IP=$(ifconfig "${DEFAULT_ETH}" | grep inet\  | awk '{print $2}')

if [[ "${DEBUG}" =~ ^[Tt]rue$ ]]; then
    echo "[1;33m *** "
    echo "REBOOT CANCELLED BECAUSE DEBUG IS ENABLED"
    echo ""
    echo "FOR REMOTE ACCESS, CURRENT IP ADDRESS IS: ${CURRENT_IP}"
    echo " *** [0m"
    killall -15 tee
    exit 0
else
    killall -15 tee
    reboot
fi
