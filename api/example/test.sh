#/bin/bash

set -eu

unset CLICOLOR
unset GREP_OPTIONS

URL="http://localhost:8080"

if [ -n "${DOCKER_HOST:-}" ]; then
  DOCKER_IP_PORT=${DOCKER_HOST#tcp://}
  DOCKER_IP=${DOCKER_IP_PORT%:*}
  URL="http://${DOCKER_IP}:8080"
fi

APIURL="${URL}/api/v1"

cpt() {
    printf "\n\n########## $1 ##########\n\n"
    printf "$2: $3\n\n"
}

fatal() {
    echo "$1" >&2
    exit 1
}

which curl &>/dev/null || fatal "Missing binary: curl"

url="${APIURL}/hosts/test1.example.com"
cpt "Create host 1" "PUT" "${url}"
curl -s -i -H "Content-Type: application/json" -X PUT -d "$(cat test1.example.com.json)" "${url}"

url="${APIURL}/hosts/test2.example.com"
cpt "Create host 2" "PUT" "${url}"
curl -s -i -H "Content-Type: application/json" -X PUT -d "$(cat test2.example.com.json)" "${url}"

url="${APIURL}/hosts/test1.example.com"
cpt "Update host 1" "PUT" "${url}"
curl -s -i -H "Content-Type: application/json" -X PUT -d "$(cat test1.example.com-update.json)" "${url}"

url="${APIURL}/hosts/test1.example.com"
cpt "Patch host 1" "PATCH" "${url}"
curl -s -i -H "Content-Type: application/json" -X PATCH -d "$(cat test1.example.com-patch.json)" "${url}"

url="${APIURL}/hosts/test1.example.com"
cpt "Delete host 1" "DELETE" "${url}"
curl -s -i -H "Content-Type: application/json" -X DELETE "${url}"

url="${APIURL}/hosts"
cpt "Get hosts" "DELETE" "${url}"
curl -s -i -H "Content-Type: application/json" "${url}"

echo
