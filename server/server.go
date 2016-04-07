package main

/*
This file initializes the
*/

import (
	// General
	"flag"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"github.com/golang/glog"
	_ "github.com/gorilla/mux"
	"github.com/iambc/xerrors"

	//DB
	"database/sql"

	_ "github.com/lib/pq"
)

var api abc_api

func main() {
	flag.Parse()
	var err error
	dbConnString = os.Getenv("ABC_DB_CONN_STRING")
	dbh, err = sql.Open("postgres", dbConnString)
	if err != nil {
		glog.Fatal("Connection to the database has failed")
	}

	wrdb := writerrdb{dbh}
	api.wr = &wrdb

	dbConnString = os.Getenv("ABC_DB_CONN_STRING")

	go http.ListenAndServe(":"+os.Getenv("ABC_FILES_SERVER_URL"), http.FileServer(http.Dir(os.Getenv("ABC_FILES_DIR"))))
	http.HandleFunc("/api", QueryStringHandler)
	http.ListenAndServe(`:`+os.Getenv("ABC_SERVER_ENDPOINT_URL"), nil)
}

func QueryStringHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	values := req.URL.Query()
	command, isPassed := values[`command`]
	if !isPassed {
		res.Write([]byte(`{"Status":"error","Msg":"Parameter 'command' is undefined.","Payload":null}`))
		return
	}

	var resp []byte
	var err error
	if command[0] == `getBoards` {
		apiKey, isPassed := values[`api_key`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'api_key' is undefined.","Payload":null}`))
			return
		}
		resp, err = api.getBoards(apiKey[0])

	} else if command[0] == `getActiveThreadsForBoard` {
		apiKey, isPassed := values[`api_key`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'api_key' is undefined.","Payload":null}`))
			return
		}
		boardID, isPassed := values[`board_id`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'board_id' is undefined.","Payload":null}`))
			return
		}
		boardIDInt, err := strconv.Atoi(boardID[0])
		if err != nil {
			res.Write([]byte(`{"Status":"error","Msg":"Wrong value for parameter board_id","Payload":null}`))
			return
		}
		resp, err = api.getActiveThreadsForBoard(apiKey[0], boardIDInt)

	} else if command[0] == `getPostsForThread` {
		apiKey, isPassed := values[`api_key`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'api_key' is undefined.","Payload":null}`))
			return
		}
		threadID, isPassed := values[`thread_id`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'thread_id' is undefined.","Payload":null}`))
			return
		}
		threadIDInt, err := strconv.Atoi(threadID[0])
		if err != nil {
			res.Write([]byte(`{"Status":"error","Msg":"Wrong value for parameter board_id","Payload":null}`))
			return
		}
		resp, err = api.getPostsForThread(apiKey[0], threadIDInt)

	} else if command[0] == `addPostToThread` {
		threadID, isPassed := values[`thread_id`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'board_id' is undefined.","Payload":null}`))
			return
		}
		threadIDInt, err := strconv.Atoi(threadID[0])
		if err != nil {
			res.Write([]byte(`{"Status":"error","Msg":"Wrong value for parameter board_id","Payload":null}`))
			return
		}
		threadBodyPost, isPassed := values[`thread_body_post`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'thread_body_post' is undefined.","Payload":null}`))
			return
		}
		attachmentUrl, isPassed := values[`attachment_id`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'attachment_id' is undefined.","Payload":null}`))
			return
		}
		clientRemoteAddr, isPassed := values[`clientRemoteAddr`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'clientRemoteAddr' is undefined.","Payload":null}`))
			return
		}
		resp, err = api.addPostToThread(threadIDInt, threadBodyPost[0], &attachmentUrl[0], clientRemoteAddr[0])

	} else if command[0] == `addThread` {
		boardID, isPassed := values[`board_id`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'board_id' is undefined.","Payload":null}`))
			return
		}
		boardIDInt, err := strconv.Atoi(boardID[0])
		if err != nil {
			res.Write([]byte(`{"Status":"error","Msg":"Wrong value for parameter board_id","Payload":null}`))
			return
		}
		threadName, isPassed := values[`thread_name`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'thread_name' is undefined.","Payload":null}`))
			return
		}
		resp, err = api.addThread(boardIDInt, threadName[0])

	} else {
		res.Write([]byte(`{"Status":"error","Msg":"No such command exists.","Payload":null}`))
		glog.Error("command: ", command[0])
	}

	if err != nil {
		if string(reflect.TypeOf(err).Name()) == `XError` {
			res.Write([]byte(`{"Status":"` + err.(xerrors.XError).Code + `","Msg":"` + err.Error() + `","Payload":null}`))
		} else {
			res.Write([]byte(`{"Status":"000","Msg":"Application Error!","Payload":null}`))
		}
		glog.Error(err)
		return
	}
	res.Write(resp)

}
