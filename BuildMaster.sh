#!/bin/bash
dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
export PATH=/home/yimin/tool/go/bin:$PATH
export GOPATH=$dir
rm -rf $dir/bin_master
mkdir $dir/bin_master
mkdir $dir/bin_master/log
cd $dir/src/xqdfs/master
go build -v
cp -rf $dir/src/xqdfs/master/master $dir/bin_master
cp -rf $dir/src/xqdfs/master/store.toml $dir/bin_master
cp -rf $dir/src/xqdfs/master/Master.sh $dir/bin_master
cp -rf $dir/webroot $dir/bin_master
chmod -R 777 $dir/bin_master
chmod -R 777 $dir/bin_master/*


