#!/bin/sh

export ABC_SERVER_ENDPOINT_URL="8089"
export ABC_DB_CONN_STRING="user=abc_api password=123 dbname=abc_dev_cluster sslmode=disable" 
# Run file for the project
/usr/local/go/bin/go run server/server.go -logtostderr=true
