#!/bin/bash
PATH_TO_DISTRIB=../../../bin/dcore
echo building 'config' ...
go build -o $PATH_TO_DISTRIB/config.exe ../cmd/config/
echo successfully
$PATH_TO_DISTRIB/config.exe
echo building 'client' ...
go build -o $PATH_TO_DISTRIB/client.exe ../cmd/client/
echo successfully
echo building 'point' ...
go build -o $PATH_TO_DISTRIB/point.exe ../cmd/point/
echo successfully
echo building 'node' ...
go build -o $PATH_TO_DISTRIB/node.exe ../cmd/node/
echo successfully
# echo generating base config...
# config can create cfg file with cmd but not with mingw. WTF !?
# $PATH_TO_DISTRIB/config.exe
# echo success