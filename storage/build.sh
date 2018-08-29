#!/bin/bash
dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
export PATH=$dir/../../../../go/bin:$PATH
export GOPATH=$dir/../../../../FaceGo
rm -rf storage
go build -v

