#!/bin/sh
. ./functions.sh
echorun go build -v "$@" -o ./bin/ ./...
