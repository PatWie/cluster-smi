#!/bin/sh
docker run --rm -it -v $(pwd):/gopath/src/github.com/patwie/cluster-smi cluster-smi-docker make
