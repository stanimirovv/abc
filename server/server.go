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

var api abcAPI

func main() {
	flag.Parse()
	var (
		dbConnString       = os.Getenv("ABC_DB_CONN_STRING")
		filesServerAddr    = os.Getenv("ABC_FILES_SERVER_URL")
		filesDir           = os.Getenv("ABC_FILES_DIR")
		serverEndpointAddr = os.Getenv("ABC_SERVER_ENDPOINT_URL")
	)

	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		glog.Fatal("Connection to the database has failed")
	}

	api.wr = &writerrdb{db}

	go http.ListenAndServe(":"+filesServerAddr, http.FileServer(http.Dir(filesDir)))
	http.HandleFunc("/api", QueryStringHandler)
	http.ListenAndServe(`:`+serverEndpointAddr, nil)
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
		apiKey, isPassed := values["api_key"]
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
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'board_Id' is undefined.","Payload":null}`))
			return
		}
		boardIDInt, err := strconv.Atoi(boardID[0])
		if err != nil {
			res.Write([]byte(`{"Status":"error","Msg":"Wrong value for parameter board_Id","Payload":null}`))
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
			res.Write([]byte(`{"Status":"error","Msg":"Wrong value for parameter board_Id","Payload":null}`))
			return
		}
		resp, err = api.getPostsForThread(apiKey[0], threadIDInt)

	} else if command[0] == `addPostToThread` {
		threadID, isPassed := values[`thread_id`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'thread_id' is undefined.","Payload":null}`))
			return
		}
		threadIDInt, err := strconv.Atoi(threadID[0])
		if err != nil {
			res.Write([]byte(`{"Status":"error","Msg":"Wrong value for parameter board_Id","Payload":null}`))
			return
		}
		threadBodyPost, isPassed := values[`thread_post_body`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'thread_post_body' is undefined.","Payload":null}`))
			return
		}
		attachmentURL, isPassed := values[`attachment_url`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'attachment_url' is undefined.","Payload":null}`))
			return
		}

		/*
			clientRemoteAddr, isPassed := values[`clientRemoteAddr`]
			if !isPassed {
				res.Write([]byte(`{"Status":"error","Msg":"Parameter 'clientRemoteAddr' is undefined.","Payload":null}`))
				return
			}
		*/
		resp, err = api.addPostToThread(threadIDInt, threadBodyPost[0], &attachmentURL[0], "127.0.0.1")
	} else if command[0] == `addThread` {
		boardID, isPassed := values[`board_id`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'board_id' is undefined.","Payload":null}`))
			return
		}
		boardIDInt, err := strconv.Atoi(boardID[0])
		if err != nil {
			res.Write([]byte(`{"Status":"error","Msg":"Wrong value for parameter board_Id","Payload":null}`))
			return
		}
		threadName, isPassed := values[`thread_name`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'thread_name' is undefined.","Payload":null}`))
			return
		}
		resp, err = api.addThread(boardIDInt, threadName[0])

	} else if command[0] == `getImageBoardClusterByApiKey` {
		apiKey, isPassed := values["api_key"]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'api_key' is undefined.","Payload":null}`))
			return
		}
		resp, err = api.getImageBoardClusterByApiKey(apiKey[0])

	} else {
		res.Write([]byte(`{"Status":"error","Msg":"No such command exists.","Payload":null}`))
		glog.Error("command: ", command[0])
	}
	/*else if command[0] == `getConfig` {
		boardID, isPassed := values[`board_id`]
		if !isPassed {
			res.Write([]byte(`{"Status":"error","Msg":"Parameter 'board_id' is undefined.","Payload":null}`))
			return
		}
		boardIDInt, err := strconv.Atoi(boardID[0])
		resp, err := api.getConfig(boardIDInt)

	}*/

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
