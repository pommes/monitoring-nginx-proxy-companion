#!/usr/bin/env bash

go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/main ./src