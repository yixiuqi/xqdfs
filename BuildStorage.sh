#!/bin/bash
dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
export PATH=/home/yimin/tool/go/bin:$PATH
export GOPATH=$dir
mkdir $dir/bin_storage
rm -rf $dir/bin_storage/*
mkdir $dir/bin_storage/log
mkdir $dir/bin_storage/binlog
cd $dir/src/xqdfs/storage
go build -v
cp -rf $dir/src/xqdfs/storage/storage $dir/bin_storage
cp -rf $dir/src/xqdfs/storage/xqdfs_storage.toml $dir/bin_storage
cp -rf $dir/src/xqdfs/storage/Storage.sh $dir/bin_storage
chmod -R 777 $dir/bin_storage
chmod -R 777 $dir/bin_storage/*


