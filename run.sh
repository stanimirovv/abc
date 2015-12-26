#!/bin/sh

export ABC_CDN_PORTNUM="6540"
export ABC_CDN_DIR="."

# Run file for the project
/usr/local/go/bin/go run cdn/cdn.go cdn/xerrors.go -logtostderr=true
