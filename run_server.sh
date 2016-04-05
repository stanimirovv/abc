#!/bin/sh

export ABC_SERVER_ENDPOINT_URL="8089"
export ABC_DB_CONN_STRING="user=abc_api password=123 dbname=abc_dev_cluster sslmode=disable" 
export ABC_FILES_DIR="./client/"
export ABC_FILES_SERVER_URL="8088"

# Run file for the project
./server/server -logtostderr=true
