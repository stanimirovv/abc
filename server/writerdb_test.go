package main

import (
	"database/sql"
	"os"
	"testing"

	"github.com/golang/glog"
)

func TestGetBoardsDbOk(t *testing.T) {
	dbConnString := os.Getenv("ABC_DB_CONN_STRING")
	db, _ := sql.Open("postgres", dbConnString)
	var api abcAPI
	wrdb := writerrdb{db}
	api.wr = &wrdb
	_, err := wrdb.getBoards(`d3c3f756aff00db5cb063765b828e87b`)
	if err != nil {
		glog.Error("Error: ", err)
		t.Fail()
	}
}

func TestGetBoardsDbFail(t *testing.T) {
	dbConnString := os.Getenv("ABC_DB_CONN_STRING")
	db, _ := sql.Open("postgres", dbConnString)
	var api abcAPI
	wrdb := writerrdb{db}
	api.wr = &wrdb
	boards, err := wrdb.getBoards(`non_existing_image_board_api_key`)
	if err != nil || len(boards) > 0 {
		t.Fail()
	}
}
