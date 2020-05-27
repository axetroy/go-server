#!/usr/bin/env bash
set -e

downloadFolder=$PWD
version=$1

targets=(
    admin
    user
    resource
    message_queue
)

get_arch() {
    a=$(uname -m)
    case ${a} in
    "x86_64" | "amd64" )
        echo "amd64"
        ;;
    "i386" | "i486" | "i586")
        echo "386"
        ;;
    *)
        echo ${NIL}
        ;;
    esac
}

get_os(){
    echo $(uname -s | awk '{print tolower($0)}')
}

main() {
    local os=$(get_os)
    local arch=$(get_arch)

    echo "Download..."
    for target in "${targets[@]}"
    do
        local dest_file="${downloadFolder}/${target}_${os}_${arch}.tar.gz"
        local asset_uri="https://github.com/axetroy/go-server/releases/download/${version}/${target}_server_${os}_${arch}.tar.gz"
        echo "Downloading ${asset_uri}"

        rm -f ${dest_file}
        curl --location --output "${dest_file}" "${asset_uri}"
    done

    exit 0
}

main