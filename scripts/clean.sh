#!/usr/bin/env bash

SOURCE="${BASH_SOURCE[0]}"
PROJECT_DIR=$(dirname $(dirname "$SOURCE") )

function walk()
{
    for file in `ls $1`
    do
        local path=$1"/"$file
        if [[ -d ${path} ]];then
            if [[ ${file} == "upload" ]] || [[ ${file} == "keys" ]]; then
                rm -rf ${path}
            else
                walk ${path}
            fi
        fi
    done
}

rm ./coverage.out

if [[ $# -ne 1 ]]
then
    walk ${PROJECT_DIR}
else
    walk $1
fi
