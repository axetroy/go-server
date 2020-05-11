#!/bin/bash

# Reference:
# https://github.com/golang/go/blob/master/src/go/build/syslist.go
os_archs=(
    darwin/386
    darwin/amd64
    linux/386
    linux/amd64
    linux/arm
    linux/arm64
    windows/386
    windows/amd64
    windows/arm
)

releases=()
targets=(
    user
    admin
    resource
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

        CGO_ENABLED=0 GOOS=${goos} GOARCH=${goarch} go build -mod=vendor -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o ./bin/${filename} ./cmd/${target}/main.go

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

