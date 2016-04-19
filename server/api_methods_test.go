package main

import (
	"testing"

	"github.com/golang/glog"
)

func TestGetBoardsOk(t *testing.T) {
	glog.Info(`TestGetBoardsOk`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.getBoards(`ok`)
	if err != nil {
		t.Fail()
	}
}

func TestGetBoardsErr(t *testing.T) {
	glog.Info(`TestGetBoardsErr`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.getBoards(`err`)
	if err == nil {
		t.Fail()
	}
}

func TestGetActiveThreadsForBoardOk(t *testing.T) {
	glog.Info(`TestGetActiveThreadsForBoardOk`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.getActiveThreadsForBoard(`ok`, 1)
	if err != nil {
		t.Fail()
	}
}

func TestGetActiveThreadsForBoardErr(t *testing.T) {
	glog.Info(`TestGetActiveThreadsForBoardErr`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.getActiveThreadsForBoard(`err`, 1)
	if err == nil {
		t.Fail()
	}
}

func TestGetActiveThreadsForBoardEmpty(t *testing.T) {
	glog.Info(`TestGetActiveThreadsForBoardEmpty`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.getActiveThreadsForBoard(`empty`, 1)
	if err != nil {
		t.Fail()
	}
}

func TestGetPostsForThreadOk(t *testing.T) {
	glog.Info(`TestGetPostsForThreadOk`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.getPostsForThread(`ok`, 1)
	if err != nil {
		t.Fail()
	}
}

func TestGetPostsForThreadEmpty(t *testing.T) {
	glog.Info(`TestGetPostsForThreadEmpty`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.getPostsForThread(`empty`, 1)
	if err != nil {
		t.Fail()
	}
}

func TestGetPostsForThreadFail(t *testing.T) {
	glog.Info(`TestGetPostsForThreadFail`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.getPostsForThread(`err`, 1)
	if err != nil {
		t.Fail()
	}
}

//addPostToThread(threadID int, threadBodyPost string, attachmentUrl *string, clientRemoteAddr string)
func TestAddThreadErrLimit(t *testing.T) {
	glog.Info(`TestAddThreadErrLimit`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.addThread(2, `a`)
	if err == nil {
		t.Fail()
	}
}

func TestAddThreadErr(t *testing.T) {
	glog.Info(`TestAddThreadErr`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.addThread(1, `a`)
	if err == nil {
		t.Fail()
	}
}

func TestAddThreadErrLimitOk(t *testing.T) {
	glog.Info(`TestAddThreadErrLimitOk`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.addThread(3, `a`)
	if err != nil {
		t.Fail()
	}
}

func TestAddThreadErrLimitUnhandledErr(t *testing.T) {
	glog.Info(`TestAddThreadErrLimitUnhandledErr`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.addThread(2, `a`)
	if err == nil {
		t.Fail()
	}
}

var str = `a`

//addPostToThread(threadID int, threadBodyPost string, attachmentUrl *string, clientRemoteAddr string)
func TestAddPostErrLimit(t *testing.T) {
	glog.Info(`TestAddPostErrLimit`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.addPostToThread(2, `a`, &str, `a`)
	if err == nil {
		t.Fail()
	}
}

func TestAddPostErr(t *testing.T) {
	glog.Info(`TestAddPostErr`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.addPostToThread(1, `a`, &str, `a`)
	if err == nil {
		t.Fail()
	}
}

func TestAddPostErrLimitLen(t *testing.T) {
	glog.Info(`TestAddPostErrLimitLen`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.addPostToThread(3, `a`, &str, `a`)
	if err == nil {
		t.Fail()
	}
}
func TestAddPostErrLimitUnhandledErr(t *testing.T) {
	glog.Info(`TestAddPostErrLimitUnhandledErr`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.addPostToThread(2, `a`, &str, `a`)
	if err == nil {
		t.Fail()
	}
}

func TestAddPost1(t *testing.T) {
	glog.Info(`TestAddPost1`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.addPostToThread(5, `a`, &str, `a`)
	if err == nil {
		t.Fail()
	}
}

func TestAddPostOk(t *testing.T) {
	glog.Info(`TestAddPostOk`)
	var api abcAPI
	wrdb := writermock{}
	api.wr = &wrdb
	_, err := api.addPostToThread(3, `aaa`, &str, `a`)
	if err != nil {
		t.Fail()
	}
}
