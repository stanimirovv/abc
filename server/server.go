package main

import (
	// General
	"github.com/golang/glog"
	"flag"
	"github.com/iambc/xerrors"
	"reflect"
	"os"

	//API
	"net/http"
	"encoding/json"

	//DB
	"database/sql"
	_ "github.com/lib/pq"
)

/*
TODO:
2) Add different input/output formats for the API
3) Add settings to the boards
4) Add settings to the threads
5) Improve the error handling
6) quote of the day
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
    BoardId int
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

    values := req.URL.Query()
    api_key := values[`api_key`][0]
    rows, err := dbh.Query("select b.id, b.name, b.descr from boards b join image_board_clusters ibc on ibc.id = b.image_board_cluster_id where api_key = $1;", api_key)
    if err != nil {
	return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `002`, true)
    }
    defer rows.Close()

    var curr_boards []boards
    for rows.Next() {
	var board boards
	err = rows.Scan(&board.Id, &board.Name, &board.Descr)
	if err != nil {
	    return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `003`, true)
	}
	curr_boards = append(curr_boards, board)
    }
    bytes, err1 := json.Marshal(api_request{"ok", nil, &curr_boards})
    if err1 != nil {
	return []byte{}, xerrors.NewUIErr(err1.Error(), err1.Error(), `004`, true)
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
	return []byte{}, xerrors.NewUIErr(`Invalid params: No board_id given!`, `Invalid params: No board_id given!`, `005`, true)
    }

    api_key := values[`api_key`][0]
    rows, err := dbh.Query(`select t.id, t.name from threads t 
				join boards b on b.id = t.board_id 
				join image_board_clusters ibc on ibc.id = b.image_board_cluster_id 
			    where t.is_active = TRUE and t.board_id = $1 and ibc.api_key = $2;`, board_id[0], api_key)
    if err != nil {
        return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `001`, true)
    }
    defer rows.Close()

    var active_threads []threads
    for rows.Next() {
	glog.Info("Popped new thread")
        var thread threads
        err = rows.Scan(&thread.Id, &thread.Name)
        if err != nil {
            return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `001`, true)
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
        return []byte{}, xerrors.NewUIErr(err1.Error(), err1.Error(), `001`, true)
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
        return []byte{},xerrors.NewUIErr(`Invalid params: No thread_id given!`, `Invalid params: No thread_id given!`, `006`, true)
    }

    api_key := values[`api_key`][0]
    rows, err := dbh.Query(`select tp.id, tp.body 
			    from thread_posts tp join threads t on t.id = tp.thread_id 
						 join boards b on b.id = t.board_id 
						 join image_board_clusters ibc on ibc.id = b.image_board_cluster_id 
			    where tp.thread_id = $1 and ibc.api_key = $2;`, thread_id[0], api_key)
    if err != nil {
	glog.Error(err)
        return []byte{}, xerrors.NewSysErr()
    }
    defer rows.Close()

    var curr_posts []thread_posts
    for rows.Next() {
	glog.Info("new post for thread with id: ", thread_id[0])
        var curr_post thread_posts
        err = rows.Scan(&curr_post.Id, &curr_post.Body)
        if err != nil {
            return []byte{}, xerrors.NewSysErr()
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
        return []byte{}, xerrors.NewSysErr()
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
        return []byte{}, xerrors.NewUIErr(`Invalid params: No thread_id given!`, `Invalid params: No thread_id given!`, `001`, true)
    }

    thread_body_post, is_passed := values[`thread_post_body`]
    if !is_passed {
        return []byte{}, xerrors.NewUIErr(`Invalid params: No thread_post_body given!`, `Invalid params: No thread_post_body given!`, `001`, true)
    }

    attachment_urls, is_passed := values[`attachment_url`]
    var attachment_url *string
    if !is_passed{
	attachment_url = nil
    }else{
	attachment_url = &attachment_urls[0]
    }

    var is_limit_reached bool
    err := dbh.QueryRow("select (select count(*) from thread_posts  where thread_id = $1) > max_posts_per_thread  from threads where id = $1;", thread_id[0]).Scan(&is_limit_reached)
    if err != nil {
	return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `001`, true)
    }

    if is_limit_reached {
	return []byte{}, xerrors.NewUIErr(`Thread post limit reached!`, `Thread post limit reached!`, `002`, true)
    }

    _, err = dbh.Query("INSERT INTO thread_posts(body, thread_id, attachment_url) VALUES($1, $2, $3)", thread_body_post[0], thread_id[0], attachment_url)

    if err != nil {
	glog.Error(err)
        return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `002`, true)
    }

    bytes, err1 := json.Marshal(api_request{"ok", nil, nil})
    if err1 != nil {
        return []byte{}, xerrors.NewUIErr(err1.Error(), err1.Error(), `001`, true)
    }

    return bytes, nil
}


func addThread(res http.ResponseWriter, req *http.Request) ([]byte,error) {
    if req == nil || res == nil{
        return []byte{}, xerrors.NewSysErr()
    }
    values := req.URL.Query()

    thread_name, is_passed := values[`thread_name`]
    if !is_passed {
        return []byte{}, xerrors.NewUIErr(`Invalid params: No thread_name given!`, `Invalid params: No thread_name given!`, `001`, true)
    }

    board_id, is_passed := values[`board_id`]
    if !is_passed {
        return []byte{}, xerrors.NewUIErr(`Invalid params: No board_id given!`, `Invalid params: No board_id given!`, `001`, true)
    }


    var is_limit_reached bool
    err := dbh.QueryRow("select (select count(*) from threads  where board_id = $1) > thread_setting_max_thread_count  from boards where id = $1;", board_id[0]).Scan(&is_limit_reached)
    if err != nil {
	glog.Error("COULD NOT SELECT thread_count")
	return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `001`, true)
    }
    if is_limit_reached {
	return []byte{}, xerrors.NewUIErr(`Thread limit reached!`, `Thread limit reached!`, `002`, true)
    }

    _, err = dbh.Query("INSERT INTO threads(name, board_id, limits_reached_action_id) VALUES($1, $2, 1)", thread_name[0], board_id[0])

    if err != nil {
	glog.Error("INSERT FAILED")
        return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `001`, true)
    }

    bytes, err1 := json.Marshal(api_request{"ok", nil, nil})
    if err1 != nil {
        return []byte{}, xerrors.NewUIErr(err1.Error(), err1.Error(), `001`, true)
    }

    return bytes, nil
}

var dbConnString = ``
var dbh *sql.DB

// sample usage
func main() {
    flag.Parse()

    var err error
    dbConnString = os.Getenv("ABC_DB_CONN_STRING") // DB will return error if empty string
    dbh, err = sql.Open("postgres", dbConnString)
    if err != nil {
	glog.Fatal(err)
    }
   
    commands := map[string]func(http.ResponseWriter, *http.Request) ([]byte, error){
				"getBoards": getBoards,
				"getActiveThreadsForBoard": getActiveThreadsForBoard,
				"getPostsForThread": getPostsForThread,
				"addPostToThread": addPostToThread,
				"addThread": addThread,
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
					    glog.Error("command: ", command[0])
					    return
					}

					res.Header().Set("Access-Control-Allow-Origin", "*")
					bytes, err := commands[command[0]](res, req)


					if err != nil{
					    if string(reflect.TypeOf(err).Name())  == `SysErr` {
						res.Write([]byte(`{"Status":"error","Msg":"` + err.Error()  +`","Payload":null}`))
					    } else if string(reflect.TypeOf(err).Name())  == `UIErr` {
						res.Write([]byte(`{"Status":"error","Msg":"`+ err.Error() +`","Payload":null}`))
					    } else {
						res.Write([]byte(`{"Status":"error","Msg":"Application Error!","Payload":null}`))
					    }
					    glog.Error(err)
					    return
					}

					glog.Info(string(bytes))
					res.Write(bytes)
    })

    http.ListenAndServe(`:`+ os.Getenv("ABC_SERVER_ENDPOINT_URL"), nil)
}


