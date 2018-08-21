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
go build -o $PATH_TO_DISTRIB/full.exe ../cmd/full/
echo 'full' built successfully