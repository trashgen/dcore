#!/bin/bash
PATH_TO_DISTRIB=../../../bin/dcore
echo building 'config' ...
go build -o $PATH_TO_DISTRIB/config.exe ../cmd/config/
echo successfully
echo building 'client' ...
go build -o $PATH_TO_DISTRIB/client.exe ../cmd/client/
echo successfully
echo building 'point' ...
go build -o $PATH_TO_DISTRIB/point.exe ../cmd/point/
echo successfully
echo building 'node' ...
go build -o $PATH_TO_DISTRIB/node.exe ../cmd/node/
echo successfully
echo configuring...
$PATH_TO_DISTRIB/config.exe
mv httpcmdconfig.cfg $PATH_TO_DISTRIB/httpcmdconfig.cfg
mv clientconfig.cfg $PATH_TO_DISTRIB/clientconfig.cfg
mv pointconfig.cfg $PATH_TO_DISTRIB/pointconfig.cfg
mv nodeconfig.cfg $PATH_TO_DISTRIB/nodeconfig.cfg
mv meta.cfg $PATH_TO_DISTRIB/meta.cfg
echo successfully