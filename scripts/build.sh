#!/bin/bash

# Reference:
# https://github.com/golang/go/blob/master/src/go/build/syslist.go
os_archs=(
    darwin/386
    darwin/amd64
    dragonfly/amd64
    freebsd/386
    freebsd/amd64
    freebsd/arm
    linux/386
    linux/amd64
    linux/arm
    linux/arm64
    linux/ppc64
    linux/ppc64le
    linux/mips
    linux/mipsle
    linux/mips64
    linux/mips64le
    linux/s390x
    nacl/386
    nacl/amd64p32
    nacl/arm
    netbsd/386
    netbsd/amd64
    netbsd/arm
    openbsd/386
    openbsd/amd64
    openbsd/arm
    plan9/386
    plan9/amd64
    plan9/arm
    solaris/amd64
    windows/386
    windows/amd64
)

releases=()
targets=(
    user
    admin
    message_queue
)

for target in "${targets[@]}"
do
    for os_arch in "${os_archs[@]}"
    do
        goos=${os_arch%/*}
        goarch=${os_arch#*/}

        filename=${target}_server

        if [[ ${goos} == "windows" ]];then
            filename+=.exe
        fi

        echo building ${target} ${os_arch}

        CGO_ENABLED=0 GOOS=${goos} GOARCH=${goarch} go build -ldflags "-s -w" -o ./bin/${filename} ./cmd/${target}/main.go >/dev/null 2>&1

        # if build success
        if [[ $? == 0 ]];then
            releases+=(${os_arch})
            cd ./bin

            tar -czf ${target}_server_${goos}_${goarch}.tar.gz ${filename}

            rm -rf ./${filename}

            cd ../
        fi
    done

    echo "${target} release:"

    for os_arch in "${releases[@]}"
    do
        printf "\t%s\n" "${os_arch}"
    done
    echo
done

