docker_download() {
	local reg=$1 name=$2 target=$3 header="${4}" full short checksum csize
	shift
	shift
	shift
	shift

	while [ -n "$1" ]; do
		file="$target/${1}.tar.gz"
		curl -s --max-time 300 --retry 3 --retry-max-time 0 -o $file -L -H "$header" "https://${reg}/v2/${name}/blobs/${1}"

		chksum=$(sha256sum $file | awk '{ print $1 }')
		[ "$chksum" == "${1##sha256:}" ] || fatal "Download failed, checksum mismatch for file: $file"

		shift
	done
}

docker_apply() {
	local source=$1 full
	shift

	while [ -n "$1" ]; do
		file="$source/${1}.tar.gz"
		tar --overwrite -xzf $file -h -C $ROOT \
			--exclude=./boot/grub/grubenv \
			--exclude=./boot/grub2/grubenv \
			--exclude=./sys \
			--exclude=./dev \
			--exclude=./proc \
			--exclude=./etc/fstab
		shift
	done
}

docker_pull() {
	local reg=$1 name=$2 tag=$3 target=$4 header="Accept: application/json"

	# Get token for docker hub
	if [ "$reg" == "registry.hub.docker.com" ]; then
		status=$(curl -s -w "%{http_code}" -H "${header}" -o $OUT "https://auth.docker.io/token?service=registry.docker.io&scope=repository:${name}:pull")
		token=$(jq -r .token $OUT)
		header="Authorization: Bearer $token"
	fi

	# Get layers
	status=$(curl -s -w "%{http_code}" -o $OUT -H "${header}" "https://${reg}/v2/${name}/manifests/${tag}")
	[ $status -ne 200 ] && fatal "Failed to get Docker image manifest"

	case $(jq -r .schemaVersion $OUT) in
	1) layers=$(jq -r '.fsLayers|reverse|.[].blobSum' $OUT) ;;
	2) layers=$(jq -r '.layers|reverse|.[].digest' $OUT) ;;
	*) warn "Unsupported schema version: $schemaVer" ;;
	esac

	set +eu
	docker_download $reg $name $target "${header}" $layers

	docker_apply $target $layers
	set -eu
}
