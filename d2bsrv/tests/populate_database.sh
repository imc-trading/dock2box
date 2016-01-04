#/bin/bash

set -eu

unset CLICOLOR
unset GREP_OPTIONS

TMPFILE="/tmp/test.json"
APIVERS="v1"

create() {
    local endp="$1" resp

    echo "POST: http://localhost:8080/${APIVERS}/${endp}" >&2
    echo "PAYLOAD:" >&2
    echo "$(cat $TMPFILE)" >&2

    resp=$(curl -s -H "Content-Type: application/json" -X POST -d "$(cat $TMPFILE)" "http://localhost:8080/${APIVERS}/${endp}")

    echo "RESPONSE:" >&2
    echo "$resp" >&2

    code=$(echo $resp | jq -r .code)
    if [ $code -ne 201 ]; then
        echo $resp
        exit $code
    fi

    echo $resp | jq -r .data.id
}

get() {
    local endp="$1" id="$2"

    echo "GET: http://localhost:8080/${APIVERS}/${endp}/id/${id}"
    echo "DATA:"
    curl -s -H "Content-Type: application/json" "http://localhost:8080/${APIVERS}/${endp}/id/${id}"
}

cpt() {
    printf "\n########## $1 ##########\n\n"
}

fatal() {
    local msg="$1"

    echo "$msg" >&2
    exit 1
}

# Check for pre. requisites
which jq &>/dev/null || fatal "Missing pre. requisite: jq"
which curl &>/dev/null || fatal "Missing pre. requisite: curl"

#
# Tenant
#
cat << EOF > $TMPFILE
{
    "tenant": "test1"
}
EOF

cpt "Create Tenant"
tenant_id=$(create "tenants")
cpt "Get Tenant"
get "tenants" $tenant_id

#
# Site
#
cat << EOF > $TMPFILE
{
    "site": "test1",
    "domain": "example.com",
    "dns": [ "192.168.0.252", "192.168.0.253" ],
    "dockerRegistry": "registry.example.com",
    "artifactRepository": "repository.example.com",
    "namingScheme": "serial-number",
    "pxeTheme": "night",
    "subnets": [
        {
            "subnet": "192.168.0.0/24",
            "mask": "255.255.255.0",
            "gw": "192.168.0.254"
        }
    ]
}
EOF

cpt "Create Site"
site_id=$(create "sites")
cpt "Get Site"
get "sites" $site_id

#
# Subnet
#
cat << EOF > $TMPFILE
{
    "subnet": "192.168.0.0/24",
    "mask": "255.255.255.0",
    "gw": "192.168.0.254",
    "siteId": "${site_id}"
}
EOF

cpt "Create Subnet"
subnet_id=$(create "subnets")
cpt "Get Subnet"
get "subnets" $subnet_id

#
# Boot Image
#
cat << EOF > $TMPFILE
{
    "image": "test1",
    "kOpts": ""
}
EOF

cpt "Create Boot Image"
boot_image_id=$(create "boot-images")
cpt "Get Boot Image"
get "boot-images" $boot_image_id

#
# Boot Image Version
#
cat << EOF > $TMPFILE
{
  "version": "latest",
  "created": "2006-01-02T15:04:05Z",
  "bootImageId": "${boot_image_id}"
}
EOF

cpt "Create Boot Image Version"
boot_image_version_id=$(create "boot-image-versions")
cpt "Get Boot Image Version"
get "boot-image-versions" $boot_image_version_id

#
# Image 
#
cat << EOF > $TMPFILE
{
    "image": "test1",
    "type": "docker",
    "bootImageId": "${boot_image_id}",
    "bootImageVersion": "latest",
    "versions": [
        {
            "version": "latest",
            "created": "2006-01-02T15:04:05Z"
        }
    ]
}
EOF

cpt "Create Image"
image_id=$(create "images")
cpt "Get Image"
get "images" $image_id

#
# Image Version
#
cat << EOF > $TMPFILE
{
  "version": "latest",
  "created": "2006-01-02T15:04:05Z",
  "imageId": "${image_id}"
}
EOF

cpt "Create Image Version"
image_version_id=$(create "image-versions")
cpt "Get Image Version"
get "image-versions" $image_version_id

#
# Host
#
cat << EOF > $TMPFILE
{
    "host": "test1.example.com",
    "build": true,
    "debug": true,
    "gpt": false,
    "imageId": "${image_id}",
    "version": "latest",
    "kOpts": "None",
    "tenantId": "${tenant_id}",
    "labels": [
        "web-server"
    ],
    "siteId": "${site_id}",
    "interfaces": [
        {
            "interface": "eth0",
            "dhcp": false,
            "ipv4": "192.168.0.1",
            "hwAddr": "a1:4c:6f:31:6c:d2",
            "subnetId": "${subnet_id}"
        },
        {
            "interface": "eth1",
            "dhcp": true,
            "hwAddr": "a1:4c:6f:31:6c:d2"
        }
    ]
}
EOF

cpt "Create Host"
host_id=$(create "hosts" "$(cat $TMPFILE)")
cpt "Get Host"
get "hosts" $host_id

rm -f $TMPFILE
