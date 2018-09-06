#!/bin/bash
dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
export PATH=/home/yimin/tool/go/bin:$PATH
export GOPATH=$dir
rm -rf $dir/bin_storage
mkdir $dir/bin_storage
mkdir $dir/bin_storage/log
mkdir $dir/bin_storage/binlog
cd $dir/src/storage
go build -v
cp -rf $dir/src/storage/storage $dir/bin_storage
cp -rf $dir/src/storage/store.toml $dir/bin_storage


