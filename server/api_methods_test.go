package main

import (
        "testing"
        )



func TestGetBoardsOk(t *testing.T) {
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.getBoards(`ok`)
    if err != nil {
	t.Fail()
    }
}


func TestGetBoardsErr(t *testing.T) {
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.getBoards(`err`)
    if err == nil {
	t.Fail()
    }
}

func TestGetActiveThreadsForBoardOk(t *testing.T) {
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.getActiveThreadsForBoard(`ok`, 1)
    if err != nil {
	t.Fail()
    }
}

func TestGetActiveThreadsForBoardErr(t *testing.T) {
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.getActiveThreadsForBoard(`err`, 1)
    if err == nil {
	t.Fail()
    }
}


func TestGetActiveThreadsForBoardEmpty(t *testing.T) {
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.getActiveThreadsForBoard(`empty`, 1)
    if err != nil {
	t.Fail()
    }
}

func TestGetPostsForThreadOk(t *testing.T){
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.getPostsForThread(`ok`, 1)
    if err != nil {
	t.Fail()
    }
}


func TestGetPostsForThreadEmpty(t *testing.T){
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.getPostsForThread(`empty`, 1)
    if err != nil {
	t.Fail()
    }
}

func TestGetPostsForThreadFail(t *testing.T){
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.getPostsForThread(`err`, 1)
    if err != nil {
	t.Fail()
    }
}

//addPostToThread(threadId int, threadBodyPost string, attachmentUrl *string, clientRemoteAddr string)
func TestAddThreadErrLimit(t *testing.T){
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.addThread(2, `a`)
    if err == nil {
	t.Fail()
    }
}

func TestAddThreadErr(t *testing.T){
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.addThread(1, `a`)
    if err == nil {
	t.Fail()
    }
}


func TestAddThreadErrLimitOk(t *testing.T){
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.addThread(3, `a`)
    if err != nil {
	t.Fail()
    }
}

func TestAddThreadErrLimitUnhandledErr(t *testing.T){
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.addThread(2, `a`)
    if err == nil {
	t.Fail()
    }
}

