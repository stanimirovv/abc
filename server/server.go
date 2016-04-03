package main

import (
	// General
	"github.com/golang/glog"
	"flag"
	"github.com/iambc/xerrors"
	"reflect"
	"os"
	"strconv"
	_ "github.com/gorilla/mux"
	"net/http"

	//DB
	"database/sql"
	_ "github.com/lib/pq"
)


var dbh *sql.DB
var dbConnString string

func main() {
    flag.Parse()
    var err error
    dbh, err = sql.Open("postgres", dbConnString)
    if err != nil {
	glog.Fatal("Connection to the database has failed")
    }
    dbConnString = os.Getenv("ABC_DB_CONN_STRING") 

    go http.ListenAndServe(":"+os.Getenv("ABC_FILES_SERVER_URL"), http.FileServer(http.Dir(os.Getenv("ABC_FILES_DIR"))))
    http.HandleFunc("/api", Handler)
    http.ListenAndServe(`:`+ os.Getenv("ABC_SERVER_ENDPOINT_URL"), nil)
}

func Handler(res http.ResponseWriter, req *http.Request){
    values := req.URL.Query()
    command, isPassed := values[`command`]
    if !isPassed {
	res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'command' is undefined.","Payload":null}`))
	return
    }

    var resp []byte
    var err error
    if(command[0] == `getBoards` ) {
	apiKey, isPassed := values[`api_key`]
	if !isPassed {
	    res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'api_key' is undefined.","Payload":null}`))
	    return
	}
	resp, err = getBoards(apiKey[0])

    } else if(command[0] == `getActiveThreadsForBoard`) {
	apiKey, isPassed := values[`api_key`]
	if !isPassed {
	    res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'api_key' is undefined.","Payload":null}`))
	    return
	}
	boardId, isPassed := values[`board_id`]
	if !isPassed {
	    res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'board_id' is undefined.","Payload":null}`))
	    return
	}
	boardIdInt, err := strconv.Atoi(boardId[0])
	if err != nil {
	    res.Write([]byte(`{"Status":"error","Msg":"Wrong value for parameter board_id","Payload":null}`))
	    return
	}
	resp, err = getActiveThreadsForBoard(apiKey[0], boardIdInt)

    } else if(command[0] == `getPostsForThread`) {
	apiKey, isPassed := values[`api_key`]
	if !isPassed {
	    res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'api_key' is undefined.","Payload":null}`))
	    return
	}
	threadId, isPassed := values[`thread_id`]
	if !isPassed {
	    res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'board_id' is undefined.","Payload":null}`))
	    return
	}
	threadIdInt, err := strconv.Atoi(threadId[0])
	if err != nil {
	    res.Write([]byte(`{"Status":"error","Msg":"Wrong value for parameter board_id","Payload":null}`))
	    return
	}
	resp, err = getPostsForThread(apiKey[0], threadIdInt)

    } else if(command[0] == `addPostToThread`) {
	threadId, isPassed := values[`thread_id`]
	if !isPassed {
	    res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'board_id' is undefined.","Payload":null}`))
	    return
	}
	threadIdInt, err := strconv.Atoi(threadId[0])
	if err != nil {
	    res.Write([]byte(`{"Status":"error","Msg":"Wrong value for parameter board_id","Payload":null}`))
	    return
	}
	threadBodyPost, isPassed := values[`thread_body_post`]
	if !isPassed {
	    res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'thread_body_post' is undefined.","Payload":null}`))
	    return
	}
	attachmentUrl, isPassed := values[`attachment_id`]
	if !isPassed {
	    res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'attachment_id' is undefined.","Payload":null}`))
	    return
	}
	clientRemoteAddr, isPassed := values[`clientRemoteAddr`]
	if !isPassed {
	    res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'clientRemoteAddr' is undefined.","Payload":null}`))
	    return
	}
	resp, err = addPostToThread(threadIdInt, threadBodyPost[0], &attachmentUrl[0], clientRemoteAddr[0])


    } else if(command[0] == `addThread`) {
	boardId, isPassed := values[`board_id`]
	if !isPassed {
	    res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'board_id' is undefined.","Payload":null}`))
	    return
	}
	boardIdInt, err := strconv.Atoi(boardId[0])
	if err != nil {
	    res.Write([]byte(`{"Status":"error","Msg":"Wrong value for parameter board_id","Payload":null}`))
	    return
	}
	threadName, isPassed := values[`thread_name`]
	if !isPassed {
	    res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'thread_name' is undefined.","Payload":null}`))
	    return
	}
	resp, err = addThread(boardIdInt, threadName[0])

    } else {
	res.Write([]byte(`{"Status":"error","Msg":"No such command exists.","Payload":null}`))
	glog.Error("command: ", command[0])
    }

    res.Header().Set("Access-Control-Allow-Origin", "*")
    if err != nil{
	if string(reflect.TypeOf(err).Name())  == `XError` {
	    res.Write([]byte(`{"Status":"`+ err.(xerrors.XError).Code +`","Msg":"` + err.Error()  +`","Payload":null}`))
	} else {
	    res.Write([]byte(`{"Status":"000","Msg":"Application Error!","Payload":null}`))
	}
	glog.Error(err)
	    return
    }
    res.Write(resp)

}

