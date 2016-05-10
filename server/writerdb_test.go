package main

import (
	"database/sql"
	"os"
	"testing"

	"github.com/golang/glog"
)

func TestGetBoardsDbOk(t *testing.T) {
	glog.Info(`TestGetBoardsDbOk`)
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
	glog.Info(`TestGetBoardsDbFail`)
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

func TestGetActiveThreadsForBoardDBOk(t *testing.T) {
	glog.Info(`TestGetBoardsDbOk`)
	dbConnString := os.Getenv("ABC_DB_CONN_STRING")
	db, _ := sql.Open("postgres", dbConnString)
	var api abcAPI
	wrdb := writerrdb{db}
	api.wr = &wrdb
	_, err := wrdb.getActiveThreadsForBoard(`d3c3f756aff00db5cb063765b828e87b`, 1)
	if err != nil {
		glog.Error("Error: ", err)
		t.Fail()
	}
}

func TestGetActiveThreadsForBoardDBNonExistingBoard(t *testing.T) {
	glog.Info(`TestGetBoardsDbOk`)
	dbConnString := os.Getenv("ABC_DB_CONN_STRING")
	db, _ := sql.Open("postgres", dbConnString)
	var api abcAPI
	wrdb := writerrdb{db}
	api.wr = &wrdb
	_, err := wrdb.getActiveThreadsForBoard(`d3c3f756aff00db5cb063765b828e87b`, 112)
	if err != nil {
		glog.Error("Error: ", err)
		t.Fail()
	}
}

func TestGetPostsForThreadDBOkt(t *testing.T) {
	glog.Info(`TestGetBoardsDbOk`)
	dbConnString := os.Getenv("ABC_DB_CONN_STRING")
	db, _ := sql.Open("postgres", dbConnString)
	var api abcAPI
	wrdb := writerrdb{db}
	api.wr = &wrdb
	_, err := wrdb.getPostsForThread(`d3c3f756aff00db5cb063765b828e87b`, 30)
	if err != nil {
		glog.Error("Error: ", err)
		t.Fail()
	}
}

func TestIsThreadLimitReached(t *testing.T) {
	glog.Info(`TestIsThreadLimitReached`)
	dbConnString := os.Getenv("ABC_DB_CONN_STRING")
	db, _ := sql.Open("postgres", dbConnString)
	var api abcAPI
	wrdb := writerrdb{db}
	api.wr = &wrdb
	_, err := wrdb.isThreadLimitReached(1)
	if err != nil {
		glog.Error("Error: ", err)
		t.Fail()
	}
}

func TestIsPostLimitReached(t *testing.T) {
	glog.Info(`TestIsThreadLimitReached`)
	dbConnString := os.Getenv("ABC_DB_CONN_STRING")
	db, _ := sql.Open("postgres", dbConnString)
	var api abcAPI
	wrdb := writerrdb{db}
	api.wr = &wrdb
	_, _, err := wrdb.isPostLimitReached(1)
	if err != nil {
		glog.Error("Error: ", err)
		t.Fail()
	}
}

func TestPostAndThreadInsert(t *testing.T) {
	glog.Info(`TestPostAndThreadInsert`)
	dbConnString := os.Getenv("ABC_DB_CONN_STRING")
	db, _ := sql.Open("postgres", dbConnString)
	var api abcAPI
	wrdb := writerrdb{db}
	api.wr = &wrdb
	thread, err := wrdb.addThread(1, `Test thread`)
	if err != nil {
		glog.Error("Error: ", err)
		t.Fail()
	}

	//Posts

	err = wrdb.addPostToThread(thread.ID, "a", nil, "a")

	if err != nil {
		glog.Error("Error: ", err)
		t.Fail()
	}

	stmt, err := db.Prepare(`DELETE FROM thread_posts WHERE thread_id = $1;`)
	if err != nil {
		glog.Error("Error: ", err)
		t.Fail()
	}

	_, err = stmt.Exec(thread.ID)

	if err != nil {
		glog.Error("Error: ", err)
		t.Fail()
	}

	//Clean up thread

	stmt, err = db.Prepare(`DELETE FROM THREADS WHERE id = $1;`)
	if err != nil {
		glog.Fatal(err)
		t.Fail()
	}
	_, err = stmt.Exec(thread.ID)
	if err != nil {
		glog.Fatal(err)
		t.Fail()
	}

	glog.Info(`Thread id: `, thread.ID)
}

func TestGetImageBoardClusterByApiKeyOk(t *testing.T) {
	glog.Info(`TestGetImageBoardClusterByApiKey`)
	dbConnString := os.Getenv("ABC_DB_CONN_STRING")
	db, _ := sql.Open("postgres", dbConnString)
	var api abcAPI
	wrdb := writerrdb{db}
	api.wr = &wrdb

	_, err := wrdb.getImageBoardClusterByApiKey(`d3c3f756aff00db5cb063765b828e87b`)
	if err != nil {
		t.Fail()
	}
}

func TestGetImageBoardClusterByApiKeyFail(t *testing.T) {
	glog.Info(`TestGetImageBoardClusterByApiKey`)
	dbConnString := os.Getenv("ABC_DB_CONN_STRING")
	db, _ := sql.Open("postgres", dbConnString)
	var api abcAPI
	wrdb := writerrdb{db}
	api.wr = &wrdb

	_, err := wrdb.getImageBoardClusterByApiKey(`d3c3`)
	if err == nil {
		t.Fail()
	}
}
