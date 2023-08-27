#!/usr/bin/env bash

./build.sh

docker build -t tyranus/monitoring-nginx-proxy-companion:`git describe` -t tyranus/monitoring-nginx-proxy-companion:latest .
