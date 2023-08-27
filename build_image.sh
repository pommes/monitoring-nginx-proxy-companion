#!/usr/bin/env bash

./build.sh

docker build -t registry.komm.link/base/docker/nginx-proxy-metrics:`git describe` -t registry.komm.link/base/docker/nginx-proxy-metrics:latest .
