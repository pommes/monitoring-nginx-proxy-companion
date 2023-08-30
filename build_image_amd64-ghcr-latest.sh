#!/usr/bin/env bash

./build.sh

export DOCKER_CLI_EXPERIMENTAL=enabled
docker buildx version
docker buildx create --use --name mybuilder
docker buildx inspect mybuilder --bootstrap

docker buildx build --platform linux/amd64 -t ghcr.io/pommes/nginx-proxy-metrics:latest . --load

docker push ghcr.io/pommes/nginx-proxy-metrics:latest
