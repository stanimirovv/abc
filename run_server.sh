#!/bin/sh

export ABC_CDN_PORTNUM="6540"
export ABC_CDN_DIR="/tmp/"
export ABC_CDN_ENDPOINT_URL="http://localhost:"$ABC_CDN_PORTNUM"/cdn/"

# Run file for the project
/usr/local/go/bin/go run server/server.go -logtostderr=true
