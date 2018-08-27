#!/bin/bash
PATH_TO_DISTRIB=../../../bin/dcore
echo create folders for distrib
mkdir -p $PATH_TO_DISTRIB/config
echo success
echo building 'config' ...
go build -o $PATH_TO_DISTRIB/config.exe ../cmd/config/
echo successfully
echo building 'view' ...
go build -o $PATH_TO_DISTRIB/view.exe ../cmd/view/
echo successfully
echo building 'signal' ...
go build -o $PATH_TO_DISTRIB/signal.exe ../cmd/signal/
echo successfully
echo building 'node' ...
go build -o $PATH_TO_DISTRIB/node.exe ../cmd/node/
echo successfully
# echo generating base config...
# config can create cfg file with cmd but not with mingw. WTF !?
# $PATH_TO_DISTRIB/config.exe
# echo success