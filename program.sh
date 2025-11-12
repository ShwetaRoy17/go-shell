#!/bin/sh

set -e


(
    cd "$(dirname "$0")"
    go build -o /tmp/go-shell-build app/*.go
)

exec /tmp/go-shell-build "$@"