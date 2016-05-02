package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/glog"
)

type MyHandler struct {
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	QueryStringHandler(w, r)
}

func TestAPIMissingCommand(t *testing.T) {
	glog.Info(`TestAPI`)
	h := MyHandler{}
	server := httptest.NewServer(&h)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		glog.Fatal("err: ", err)
		t.Fail()
	}
	glog.Info(`resp: `, resp)
}

func TestAPIUnknownCommand(t *testing.T) {
	glog.Info(`TestAPIUnknownCommand`)
	h := MyHandler{}
	server := httptest.NewServer(&h)
	defer server.Close()

	resp, err := http.Get(server.URL + `?command=oewihwoweoriwuori`)
	if err != nil {
		glog.Fatal("err: ", err)
		t.Fail()
	}
	glog.Info(`resp: `, resp)
}

func TestAPIGetBoards(t *testing.T) {
	//glog.Info(`TestAPIGetBoards`)
	h := MyHandler{}
	server := httptest.NewServer(&h)
	defer server.Close()

	resp, err := http.Get(server.URL + "?command=getBoards")
	if err != nil {
		glog.Fatal("err: ", err)
		t.Fail()
	}
	glog.Info(`resp: `, resp)
}

func TestAPIGetBoardsWrongApiKey(t *testing.T) {
	//glog.Info(`TestAPIGetBoards`)
	wr := writermock{}
	api.wr = &wr

	h := MyHandler{}
	server := httptest.NewServer(&h)
	defer server.Close()

	resp, err := http.Get(server.URL + "?command=getBoards&api_key=d3c3f756aff00db5cb063765b828e87b")
	if err != nil {
		glog.Fatal("err: ", err)
		t.Fail()
	}
	glog.Info(`resp: `, resp)
}

func TestAPIGetPostsForThreadNoApiKey(t *testing.T) {
	//glog.Info(`TestAPIGetBoards`)
	wr := writermock{}
	api.wr = &wr

	h := MyHandler{}
	server := httptest.NewServer(&h)
	defer server.Close()

	resp, err := http.Get(server.URL + "?command=getPostsForThread")
	if err != nil {
		glog.Fatal("err: ", err)
		t.Fail()
	}
	glog.Info(`resp: `, resp)
}

func TestAPIGetPostsForThread(t *testing.T) {
	//glog.Info(`TestAPIGetBoards`)
	wr := writermock{}
	api.wr = &wr

	h := MyHandler{}
	server := httptest.NewServer(&h)
	defer server.Close()

	resp, err := http.Get(server.URL + "?command=getPostsForThread&api_key=d3c3f756aff00db5cb063765b828e87b&thread_id=1")
	if err != nil {
		glog.Fatal("err: ", err)
		t.Fail()
	}
	glog.Info(`resp: `, resp)
}

func TestAPIGetActiveThreadsForBoardNoApiKey(t *testing.T) {
	//glog.Info(`TestAPIGetBoards`)
	wr := writermock{}
	api.wr = &wr

	h := MyHandler{}
	server := httptest.NewServer(&h)
	defer server.Close()

	resp, err := http.Get(server.URL + "?command=getActiveThreadsForBoard")
	if err != nil {
		glog.Fatal("err: ", err)
		t.Fail()
	}
	glog.Info(`resp: `, resp)
}

func TestAPIGetActiveThreadsForBoardMissingBoardId(t *testing.T) {
	//glog.Info(`TestAPIGetBoards`)
	wr := writermock{}
	api.wr = &wr

	h := MyHandler{}
	server := httptest.NewServer(&h)
	defer server.Close()

	resp, err := http.Get(server.URL + "?command=getActiveThreadsForBoard&api_key=d3c3f756aff00db5cb063765b828e87b&thread_id=1")
	if err != nil {
		glog.Fatal("err: ", err)
		t.Fail()
	}
	glog.Info(`resp: `, resp)
}

func TestAPIGetActiveThreadsForBoardMissingThread(t *testing.T) {
	//glog.Info(`TestAPIGetBoards`)
	wr := writermock{}
	api.wr = &wr

	h := MyHandler{}
	server := httptest.NewServer(&h)
	defer server.Close()

	resp, err := http.Get(server.URL + "?command=getActiveThreadsForBoard&api_key=d3c3f756aff00db5cb063765b828e87b&thread_id1=1")
	if err != nil {
		glog.Fatal("err: ", err)
		t.Fail()
	}
	glog.Info(`resp: `, resp)
}

func TestAPIGetActiveThreadsForBoardWrongThreadId(t *testing.T) {
	//glog.Info(`TestAPIGetBoards`)
	wr := writermock{}
	api.wr = &wr

	h := MyHandler{}
	server := httptest.NewServer(&h)
	defer server.Close()

	resp, err := http.Get(server.URL + "?command=getActiveThreadsForBoard&api_key=d3c3f756aff00db5cb063765b828e87b&thread_id=a")
	if err != nil {
		glog.Fatal("err: ", err)
		t.Fail()
	}
	glog.Info(`resp: `, resp)
}

func TestAPIGetActiveThreadsForBoardStringBoardId(t *testing.T) {
	//glog.Info(`TestAPIGetBoards`)
	wr := writermock{}
	api.wr = &wr

	h := MyHandler{}
	server := httptest.NewServer(&h)
	defer server.Close()

	resp, err := http.Get(server.URL + "?command=getActiveThreadsForBoard&api_key=d3c3f756aff00db5cb063765b828e87b&board_id=asd")
	if err != nil {
		glog.Fatal("err: ", err)
		t.Fail()
	}
	glog.Info(`resp: `, resp)
}

func TestAPIGetActiveThreadsForBoard(t *testing.T) {
	//glog.Info(`TestAPIGetBoards`)
	wr := writermock{}
	api.wr = &wr

	h := MyHandler{}
	server := httptest.NewServer(&h)
	defer server.Close()

	resp, err := http.Get(server.URL + "?command=getActiveThreadsForBoard&api_key=d3c3f756aff00db5cb063765b828e87b&board_id=1")
	if err != nil {
		glog.Fatal("err: ", err)
		t.Fail()
	}
	glog.Info(`resp: `, resp)
}
