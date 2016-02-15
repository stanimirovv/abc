package main

import (
	// General
	"github.com/golang/glog"
	"flag"
	"github.com/iambc/xerrors"
	"reflect"

	//API
	"net/http"
	"encoding/json"

	//DB
	"database/sql"
	_ "github.com/lib/pq"
)

/*
TODO:
ADD THREAD function
1) DB connecting must not be hardcoded
2) Add different input/output formats for the API
3) Start adding settings to the boards, etc
*/


type image_board_clusters struct {
    Id int
    Descr string
    LongDescr string
    BoardLimitCount int
}

type boards struct {
    Id int
    Name string
    Descr string
    ImageBoardClusterId string
    MaxThreadCount int //to be checked in insert thread
    MaxActiveThreadCount int //to be checked  in insert thread
    MaxPostsPerThread int // to be checked in insert thread
    AreAttachmentsAllowed bool // to be checked in insert post
    PostLimitsReachedActionId int // to be checked in insert post
}

type threads struct{
    Id int
    Name string
    Descr string
    Board_id int
    MaxPostsPerThread int
    AreAttachmentsAllowed bool
    LimitsReachedActionId int
}

type thread_posts struct{
    Id int
    Body string
    ThreadId int
    AttachmentUrl int
}

type thread_limits_reached_actions struct{
    Id	    int
    Name    string
    Descr   string
}

type api_request struct{
    Status  string
    Msg	    *string
    Payload interface{}
}


func getBoards(res http.ResponseWriter, req *http.Request)  ([]byte, error) {
    if req == nil || res == nil {
	return []byte{}, xerrors.NewSysErr()
    }

    dbh, err := sql.Open("postgres", "user=abc_api password=123 dbname=abc_dev_cluster sslmode=disable")
    if err != nil {
	return []byte{}, xerrors.NewUiErr(err.Error(), err.Error())
    }

    values := req.URL.Query()
    api_key := values[`api_key`][0]
    rows, err := dbh.Query("select b.id, b.name, b.descr from boards b join image_board_clusters ibc on ibc.id = b.id where api_key = $1;", api_key)
    if err != nil {
	return []byte{}, xerrors.NewUiErr(err.Error(), err.Error())
    }
    defer rows.Close()

    var curr_boards []boards
    for rows.Next() {
	var board boards
	err = rows.Scan(&board.Id, &board.Name, &board.Descr)
	if err != nil {
	    return []byte{}, xerrors.NewUiErr(err.Error(), err.Error())
	}
	curr_boards = append(curr_boards, board)
    }
    bytes, err1 := json.Marshal(api_request{"ok", nil, &curr_boards})
    if err1 != nil {
	return []byte{}, xerrors.NewUiErr(err1.Error(), err1.Error())
    }
    return bytes, nil
}


func getActiveThreadsForBoard(res http.ResponseWriter, req *http.Request)  ([]byte, error) {
    if req == nil || res == nil {
        return []byte{}, xerrors.NewSysErr()
    }
    values := req.URL.Query()
    board_id, is_passed := values[`board_id`]
    if !is_passed {
	return []byte{}, xerrors.NewUiErr(`Invalid params: No board_id given!`, `Invalid params: No board_id given!`)
    }

    dbh, err := sql.Open("postgres", "user=abc_api password=123 dbname=abc_dev_cluster sslmode=disable")
    if err != nil {
        return []byte{}, xerrors.NewUiErr(err.Error(), err.Error())
    }

    api_key := values[`api_key`][0]
    rows, err := dbh.Query("select t.id, t.name from threads t 
				join boards b on b.id = t.board_id 
				join image_board_clusters ibc on ibc.id = b.id 
			    where t.is_active = TRUE and t.board_id = $1 and bc.api_key = $2;", board_id[0], api_key)
    if err != nil {
        return []byte{}, xerrors.NewUiErr(err.Error(), err.Error())
    }
    defer rows.Close()

    var active_threads []threads
    for rows.Next() {
	glog.Info("Popped new thread")
        var thread threads
        err = rows.Scan(&thread.Id, &thread.Name)
        if err != nil {
            return []byte{}, xerrors.NewUiErr(err.Error(), err.Error())
        }
        active_threads = append(active_threads, thread)
    }
    var bytes []byte
    var err1 error
    if(len(active_threads) == 0){
        errMsg := "No objects returned."
        bytes, err1 = json.Marshal(api_request{"error", &errMsg, &active_threads})
    }else {
        bytes, err1 = json.Marshal(api_request{"ok", nil, &active_threads})
    }

    if err1 != nil {
        return []byte{}, xerrors.NewUiErr(err1.Error(), err1.Error())
    }

    return bytes, nil
}


func getPostsForThread(res http.ResponseWriter, req *http.Request)  ([]byte, error) {
    if req == nil || res == nil {
        return []byte{}, xerrors.NewSysErr()
    }
    values := req.URL.Query()
    thread_id, is_passed := values[`thread_id`]
    if !is_passed {
        return []byte{},xerrors.NewUiErr(`Invalid params: No thread_id given!`, `Invalid params: No thread_id given!`)
    }

    dbh, err := sql.Open("postgres", "user=abc_api password=123 dbname=abc_dev_cluster sslmode=disable")
    if err != nil {
        return []byte{}, xerrors.NewUiErr(err.Error(), err.Error())
    }

    api_key := values[`api_key`][0]
    rows, err := dbh.Query("select tp.id, tp.body 
			    from thread_posts tp join threads t on t.id = tp.thread_id 
						 join boards b on b.id = t.board_id 
						 join image_board_clusters ibc on ibc.id = b.id 
			    where tp.thread_id = $1 and ibc.api_key = $2;", thread_id[0], api_key)
    if err != nil {
        return []byte{}, xerrors.NewUiErr(err.Error(), err.Error())
    }
    defer rows.Close()

    var curr_posts []thread_posts
    for rows.Next() {
        var curr_post thread_posts
        err = rows.Scan(&curr_post.Id, &curr_post.Body)
        if err != nil {
            return []byte{}, xerrors.NewUiErr(err.Error(), err.Error())
        }
        curr_posts = append(curr_posts, curr_post)
    }

    var bytes []byte
    var err1 error
    if(len(curr_posts) == 0){
	errMsg := "No objects returned."
	bytes, err1 = json.Marshal(api_request{"error", &errMsg, &curr_posts})
    }else {
	bytes, err1 = json.Marshal(api_request{"ok", nil, &curr_posts})
    }

    if err1 != nil {
        return []byte{}, xerrors.NewUiErr(err1.Error(), err1.Error())
    }

    return bytes, nil
}


func addPostToThread(res http.ResponseWriter, req *http.Request) ([]byte,error) {
    if req == nil || res == nil{
        return []byte{}, xerrors.NewSysErr()
    }
    values := req.URL.Query()
    thread_id, is_passed := values[`thread_id`]
    if !is_passed {
        return []byte{}, xerrors.NewUiErr(`Invalid params: No thread_id given!`, `Invalid params: No thread_id given!`)
    }

    thread_body_post, is_passed := values[`thread_post_body`]
    if !is_passed {
        return []byte{}, xerrors.NewUiErr(`Invalid params: No thread_post_body given!`, `Invalid params: No thread_post_body given!`)
    }

    attachment_urls, is_passed := values[`attachment_url`]
    var attachment_url *string
    if !is_passed{
	attachment_url = nil
    }else{
	attachment_url = &attachment_urls[0]
    }

    dbh, err := sql.Open("postgres", "user=abc_api password=123 dbname=abc_dev_cluster sslmode=disable")
    if err != nil {
        return []byte{}, xerrors.NewUiErr(err.Error(), err.Error())
    }

    _, err = dbh.Query("INSERT INTO thread_posts(body, thread_id, attachment_url) VALUES($1, $2, $3)", thread_body_post[0], thread_id[0], attachment_url)

    if err != nil {
        return []byte{}, xerrors.NewUiErr(err.Error(), err.Error())
    }

    bytes, err1 := json.Marshal(api_request{"ok", nil, nil})
    if err1 != nil {
        return []byte{}, xerrors.NewUiErr(err1.Error(), err1.Error())
    }

    return bytes, nil
}


// sample usage
func main() {
    flag.Parse()

    commands := map[string]func(http.ResponseWriter, *http.Request) ([]byte, error){
				"getBoards": getBoards,
				"getActiveThreadsForBoard": getActiveThreadsForBoard,
				"getPostsForThread": getPostsForThread,
				"addPostToThread": addPostToThread,
				//TODO add thread 
			       }

    http.HandleFunc("/api", func(res http.ResponseWriter, req *http.Request) {
					values := req.URL.Query()
					command, is_passed := values[`command`]
					if !is_passed {
					    res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'command' is undefined.","Payload":null}`))
					    return
					}


					_, is_passed = values[`api_key`]
					if !is_passed {
					    res.Write([]byte(`{"Status":"error","Msg":"Paremeter 'api_key' is undefined.","Payload":null}`))
					    return
					}

					_, is_passed = commands[command[0]]
					if !is_passed{
					    res.Write([]byte(`{"Status":"error","Msg":"No such command exists.","Payload":null}`))
					    return
					}


					bytes, err := commands[command[0]](res, req)


					if err != nil{
					    if string(reflect.TypeOf(err).Name())  == `SysErr` {
						res.Write([]byte(`{"Status":"error","Msg":"Application error!","Payload":null}`))
					    }else if string(reflect.TypeOf(err).Name())  == `UiErr` {
						res.Write([]byte(`{"Status":"error","Msg":"`+ err.Error() +`","Payload":null}`))
					    }
					    glog.Error(err)
					    return
					}

					glog.Info(string(bytes))
					res.Write(bytes)
    })

    http.ListenAndServe(`:8089`, nil)
}


