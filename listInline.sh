#!/bin/sh
go build -gcflags='-m' -o /dev/null ./... 2>&1 | grep -e 'can inline' -e 'inlining'
