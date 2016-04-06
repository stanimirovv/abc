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

var str = `a`

//addPostToThread(threadId int, threadBodyPost string, attachmentUrl *string, clientRemoteAddr string)
func TestAddPostErrLimit(t *testing.T){
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.addPostToThread(2, `a`, &str, `a`)
    if err == nil {
	t.Fail()
    }
}

func TestAddPostErr(t *testing.T){
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.addPostToThread(1, `a`, &str, `a`)
    if err == nil {
	t.Fail()
    }
}

func TestAddPostErrLimitLen(t *testing.T){
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.addPostToThread(3, `a`, &str, `a`)
    if err == nil {
	t.Fail()
    }
}
func TestAddPostErrLimitUnhandledErr(t *testing.T){
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.addPostToThread(2, `a`, &str, `a`)
    if err == nil {
	t.Fail()
    }
}

func TestAddPost1(t *testing.T){
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.addPostToThread(5, `a`, &str, `a`)
    if err == nil {
	t.Fail()
    }
}

func TestAddPostOk(t *testing.T){
    var api abc_api
    wrdb := writermock{}
    api.wr = &wrdb
    _, err := api.addPostToThread(3, `aaa`, &str, `a`)
    if err != nil {
	t.Fail()
    }
}

