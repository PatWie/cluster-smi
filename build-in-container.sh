#!/bin/sh
podman run --rm -it -v $(pwd):/gopath/src/github.com/patwie/cluster-smi cluster-smi-docker make
