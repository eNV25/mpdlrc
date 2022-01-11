#!/bin/sh
eval "$(go env | sed -E 's/(.*)=".*"/\1=; export \1/g')"
echoexec () { echo "$@"; exec "$@"; }
echoexec go build -v "$@" -o ./bin/ ./...
