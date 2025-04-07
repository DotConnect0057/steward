#!/bin/sh

errorExit() {
    echo "*** $*" 1>&2
    exit 1
}

curl -sfk --max-time 2 https://localhost:{{.port}}/healthz -o /dev/null || errorExit "Error GET https://localhost:{{.port}}/healthz"