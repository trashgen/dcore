#!/bin/bash
docker container rm --force config
docker image rm --force dcore

docker rmi --force $(docker images | grep "^<none>" | awk "{print $3}")

cd ../
docker build --rm -t dcore .
# docker run --name config dcore
winpty docker run -it --name config dcore