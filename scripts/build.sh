#!/bin/bash

# Reference:
# https://github.com/golang/go/blob/master/src/go/build/syslist.go
os_archs=(
    darwin/amd64
    linux/amd64
    windows/amd64
)

releases=()
target=$1

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
    else
        exit $?
    fi
done

echo "${target} release:"

for os_arch in "${releases[@]}"
do
    printf "\t%s\n" "${os_arch}"
done
echo
