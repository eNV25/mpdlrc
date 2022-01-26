#!/bin/sh
eval "$(go env | sed -E 's/(.*)=".*"/\1=; export \1/g')"
exec ./envbuild.sh "$@"
