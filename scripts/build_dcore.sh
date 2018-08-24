#!/bin/bash
PATH_TO_DISTRIB=../../../bin/dcore
echo create folders for distrib
mkdir -p $PATH_TO_DISTRIB/config
echo success
echo building...
go build -o $PATH_TO_DISTRIB/config.exe ../cmd/config/
echo 'config' built successfully
go build -o $PATH_TO_DISTRIB/view.exe ../cmd/view/
echo 'view' built successfully
go build -o $PATH_TO_DISTRIB/signal.exe ../cmd/signal/
echo 'signal' built successfully
go build -o $PATH_TO_DISTRIB/node.exe ../cmd/node/
echo 'node' built successfully
# echo generating base config...
# config can create cfg file with cmd but not with mingw. WTF !?
# $PATH_TO_DISTRIB/config.exe
# echo success