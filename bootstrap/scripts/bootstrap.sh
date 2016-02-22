#!/bin/bash

set -eu
#set -x

API_URL="http://dock2box:8080/api/v1"

info() {
    local message="$1"

    echo -e "\e[32m+ ${message}\e[0m"
}

warning() {
    local message="$1"

    echo -e "\e[33m! ${message}\e[0m"
}

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
    [ "${value}" == "null" -o "${value}" == "" ] && value="${default}"
    echo "${value}"
}

# Get passed paramenters from kernel options
BOOT_HWADDR=$(get_kopt "BOOTIF")
[ -n "${BOOT_HWADDR}" ] || fatal "Missing kernel option: BOOT_HWADDR"
info "Boot Hardware Address: ${BOOT_HWADDR}"

BOOT_IF=$(grep -i --with-filename ${BOOT_HWADDR} /sys/class/net/*/address | cut -d/ -f5)
[ -n "${BOOT_HWADDR}" ] || fatal "Can't determine boot interface"
info "Boot Interface: ${BOOT_IF}"

BOOT_IP=$(get_kopt "IP")
[ -n "${BOOT_IP}" ] || fatal "Missing kernel option: IP"
info "Boot IP Address: ${BOOT_IP}"

BOOT_NETMASK=$(get_kopt "NETMASK")
[ -n "${BOOT_NETMASK}" ] || fatal "Missing kernel option: NETMASK"
info "Boot Netmask: ${BOOT_NETMASK}"

BOOT_GATEWAY=$(get_kopt "GATEWAY")
[ -n "${BOOT_NETMASK}" ] || fatal "Missing kernel option: GATEWAY"
info "Boot Gateway: ${BOOT_GATEWAY}"

#BOOT_SUBNET=$(get_kopt "SUBNET")
#[ -n "${BOOT_SUBNET}" ] || fatal "Missing kernel option: SUBNET"
#info "Boot Subnet: ${BOOT_SUBNET}"

IMAGE=$(get_kopt "image_name")
[ -n "${IMAGE}" ] || fatal "Missing kernel option: image_name"
info "Image: ${IMAGE}"

IMAGE_VERSION=$(get_kopt "image_version")
[ -n "${IMAGE_VERSION}" ] || fatal "Missing kernel option: image_version"
info "Image Version: ${IMAGE_VERSION}"

IMAGE_TYPE=$(get_kopt "image_type")
[ -n "${IMAGE_TYPE}" ] || fatal "Missing kernel option: image_type"
info "Image Type: ${IMAGE_TYPE}"

BOOT_HOSTNAME=$(get_kopt "hostname")
[ -n "${BOOT_HOSTNAME}" ] || fatal "Missing kernel option: hostname"
info "Boot Hostname: ${BOOT_HOSTNAME}"

SITE=$(get_kopt "SITE")
if [ -z "${SITE}" ]; then
    warning "Missing kernel option: SITE"
    info "Will determine domain and DNS servers using DHCP instead of site"

    DOMAIN=$(awk '/domain/ {print $2}' /etc/resolv.conf | head -1)
    [ -n "${DOMAIN}" ] || fatal "Can't determine DHCP domain"

    DNS_SERVERS=$(awk '/nameserver/ {print $2}' /etc/resolv.conf | tr '\n' ' ')
    [ -n "${DNS_SERVERS}" ] || fatal "Can't determine DHCP nameserver(s)"
else
    info "Site: ${SITE}"

    # Get site configuration
    SITE_DATA=$(get_resource "/sites/${SITE}")

    DOMAIN=$(get_key "${SITE_DATA}" ".domain")
    [ -n "${DOMAIN}" ] || fatal "Failed to retrieve domain for site: ${SITE}"

    DNS_SERVERS=$(get_key "${SITE_DATA}" ".dns_servers")
    [ -n "${DNS_SERVERS}" ] || fatal "Failed to retrieve DNS servers for site: ${SITE}"
fi

info "Use Domain: ${DOMAIN}"
info "Use DNS Servers: ${DNS_SERVERS}"

# Get host configuration
FQDN="${BOOT_HOSTNAME}.${DOMAIN}"
info "Use Full Hostname: ${FQDN}"

HOST_DATA=$(get_resource "/hosts/${FQDN}")

DEBUG=$(get_key "${HOST_DATA}" ".debug")
if [[ "${DEBUG}" =~ "[Tt]rue" ]]; then
    set -x
fi

# If dhcp is not set default to false
DHCP=$(get_key_default "${HOST_DATA}" ".interface.${BOOT_IF}.dhcp" "false")

# Get network configuration
if [[ "${DHCP}" =~ ^[Ff]alse$ ]]; then
    IPV4=$(get_key "${HOST_DATA}" ".interface.${BOOT_IF}.ip")
    if [ -n "${IPV4}" ]; then
        info "Use IPv4: ${IPV4}"
    else
        warning "Failed to retrieve IPv4 address for host: ${FQDN} interface: ${BOOT_IF}"
    fi

    NETMASK=$(get_key "${HOST_DATA}" ".interface.${BOOT_IF}.netmask")
    if [ -n "${NETMASK}" ]; then
        info "Use Netmask: ${NETMASK}"
    else
        warning "Failed to retrieve netmask for host: ${FQDN} interface: ${BOOT_IF}"
    fi

    GATEWAY=$(get_key "${HOST_DATA}" ".interface.${BOOT_IF}.gateway")
    if [ -n "${GATEWAY}" ]; then
        info "Use Netmask: ${GATEWAY}"
    else
        warning "Failed to retrieve gateway address for host: ${FQDN} interface: ${BOOT_IF}"
    fi

    # If there is no network configuration default to using dhcp
    if [ -z "${IPV4}" ]; then
        warning "No network configuration defined default to DHCP"
	DHCP="true"
    fi
fi

info "Use DHCP: ${DHCP}"

# Get image configuration
IMAGE_DATA=$(get_resource "/hosts/${IMAGE}")

DISTRO=$(get_key "${IMAGE_DATA}" ".distro")
[ -n "${DISTRO}" ] || fatal "Failed to retrieve distribution for image: ${IMAGE}"

exit 0

# Get default configuartion
#SCRIPT=$(get_key "${DEFAULT}" ".install_script")

# Set initial root password to "dock2box"
/usr/sbin/usermod -p '$6$WyzcP9uG$64NvhY1P.W1d08CRxFjJ5QUDyZcrzyjxVPV82bAxBgf9eE2sdSZs.v47HGdWPvLxhNxAkrJten86R2Qb1vdhe/' root

# Print cmdline
grep BOOT --with-filename /proc/cmdline

# Print network information
ifconfig "${BOOT_IF}"

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
info "Set hostname: ${FQDN}"
hostname "${FQDN}"

info "Download install script"
download "${SCRIPT}" "/root/install.sh"

echo "Executing the install script"
chmod +x /root/install.sh
#export IP NETMASK SUBNET GATEWAY HOSTNAME HOSTNAME_SHORT NAMING_SCHEME DOMAIN BOOT_MAC BOOT_ETH DEBUG IMAGE SCRIPT BOOT_IMAGE_NAME BOOT_IMAGE_VERSION BOOT_IMAGE_TYPE DNS

#/root/install.sh

CURRENT_IP=$(ifconfig "${BOOT_IF}" | grep inet\  | awk '{print $2}')

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
