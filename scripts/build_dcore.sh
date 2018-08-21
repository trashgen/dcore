#!/bin/bash
PROJECT_NAME="dcore"
echo building...
go build -o ../../../bin/$PROJECT_NAME/config.exe ../cmd/config/
echo 'config' built successfully
go build -o ../../../bin/$PROJECT_NAME/view.exe ../cmd/view/
echo 'view' built successfully
go build -o ../../../bin/$PROJECT_NAME/signal.exe ../cmd/signal/
echo 'signal' built successfully
go build -o ../../../bin/$PROJECT_NAME/full.exe ../cmd/full/
echo 'full' built successfully
sleep 3