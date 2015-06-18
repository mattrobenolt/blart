#!/bin/bash
set -x

rm -rf bin/
docker build --rm -t blart .
docker run --rm -v $PWD:/go/src/blart -it blart
