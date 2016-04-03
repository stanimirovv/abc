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

	//API
	"net/http"
	"encoding/json"

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

    } else if(command[0] == `addPostToThread`) {

    } else if(command[0] == `addThread`) {

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

func getBoards(apiKey string)  ([]byte, error) {

    rows, err := dbh.Query("select b.id, b.name, b.descr from boards b join image_board_clusters ibc on ibc.id = b.image_board_cluster_id where api_key = $1;", apiKey)
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


func getActiveThreadsForBoard(apiKey string, boardId int)  ([]byte, error) {

    rows, err := dbh.Query(`select t.id, t.name, count(*), (select count(*) from thread_posts where thread_id = t.id and attachment_url is not null) from threads t  
				join boards b on b.id = t.board_id 
				join image_board_clusters ibc on ibc.id = b.image_board_cluster_id 
				left join thread_posts tp on tp.thread_id = t.id
			    where t.is_active = TRUE and t.board_id = $1 and ibc.api_key = $2 group by 1,2 order by t.id;`, boardId, apiKey)
    if err != nil {
        return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `006`, true)
    }
    defer rows.Close()

    var active_threads []threads
    for rows.Next() {
	glog.Info("Popped new thread")
        var thread threads
        err = rows.Scan(&thread.Id, &thread.Name, &thread.PostCount, &thread.PostCountWithAttachment)
        if err != nil {
            return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `007`, true)
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
        return []byte{}, xerrors.NewUIErr(err1.Error(), err1.Error(), `008`, true)
    }

    return bytes, nil
}


func getPostsForThread(apiKey string, threadId int)  ([]byte, error) {
    rows, err := dbh.Query(`select tp.id, tp.body, tp.attachment_url, tp.inserted_at, tp.source_ip 
			    from thread_posts tp join threads t on t.id = tp.thread_id 
						 join boards b on b.id = t.board_id 
						 join image_board_clusters ibc on ibc.id = b.image_board_cluster_id 
			    where tp.thread_id = $1 and ibc.api_key = $2 and t.is_active = true;`, threadId, apiKey)
    if err != nil {
	glog.Error(err)
        return []byte{}, xerrors.NewSysErr()
    }
    defer rows.Close()

    var curr_posts []thread_posts
    for rows.Next() {
	glog.Info("new post for thread with id: ", threadId)
        var curr_post thread_posts
        err = rows.Scan(&curr_post.Id, &curr_post.Body, &curr_post.AttachmentUrl, &curr_post.InsertedAt, &curr_post.SourceIp)
        if err != nil {
	    glog.Error(err)
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

    var is_limit_reached bool
    var max_post_length int
    var min_post_length int
    err := dbh.QueryRow("select (select count(*) from thread_posts  where thread_id = $1) > max_posts_per_thread, min_post_length, max_post_length  from threads where id = $1;", thread_id[0]).Scan(&is_limit_reached, &min_post_length, &max_post_length)
    if err != nil {
	return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `009`, true)
    }

    if is_limit_reached {
	dbh.QueryRow("UPDATE threads set is_active = false where id = $1",  thread_id[0]).Scan()
	return []byte{}, xerrors.NewUIErr(`Thread post limit reached!`, `Thread post limit reached!`, `010`, true)
    }

   if(min_post_length > len(thread_body_post[0])  && min_post_length != -1){
	return []byte{}, xerrors.NewUIErr(`Post length is less than minimum length!`, `Post length is less than minimum length! post length: ` + strconv.Itoa(len(thread_body_post[0]))  +` min length: ` + strconv.Itoa(min_post_length) , `020`, false)
    }
   if(max_post_length < len(thread_body_post[0])  && max_post_length != -1){
	return []byte{}, xerrors.NewUIErr(`Post length is more than maximum length!`, `Post length is more than maximum length! post length: ` + strconv.Itoa(len(thread_body_post[0]))  +` max length: ` + strconv.Itoa(max_post_length) , `021`, false)
    }

    attachment_urls, is_passed := values[`attachment_url`]
    var attachment_url *string
    if !is_passed{
	attachment_url = nil
    }else{
	attachment_url = &attachment_urls[0]
    }

    if(*attachment_url == ``){
	attachment_url = nil
    }

    _, err = dbh.Query("INSERT INTO thread_posts(body, thread_id, attachment_url, source_ip) VALUES($1, $2, $3, $4)", thread_body_post[0], thread_id[0], attachment_url, req.RemoteAddr)

    if err != nil {
	glog.Error(err)
        return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `011`, true)
    }

    bytes, err1 := json.Marshal(api_request{"ok", nil, nil})
    if err1 != nil {
        return []byte{}, xerrors.NewUIErr(err1.Error(), err1.Error(), `012`, true)
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
        return []byte{}, xerrors.NewUIErr(`Invalid params: No thread_name given!`, `Invalid params: No thread_name given!`, `013`, true)
    }

    board_id, is_passed := values[`board_id`]
    if !is_passed {
        return []byte{}, xerrors.NewUIErr(`Invalid params: No board_id given!`, `Invalid params: No board_id given!`, `014`, true)
    }


    var is_limit_reached bool
    err := dbh.QueryRow("select (select count(*) from threads  where board_id = $1) > thread_setting_max_thread_count  from boards where id = $1;", board_id[0]).Scan(&is_limit_reached)
    if err != nil {
	glog.Error("COULD NOT SELECT thread_count")
	return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `015`, true)
    }
    if is_limit_reached {
	return []byte{}, xerrors.NewUIErr(`Thread limit reached!`, `Thread limit reached!`, `016`, true)
    }

    var threadId int
    var threadName string
    err = dbh.QueryRow("INSERT INTO threads(name, board_id, limits_reached_action_id, max_posts_per_thread) VALUES($1, $2, 1, 10)  RETURNING id, name", thread_name[0], board_id[0]).Scan(&threadId, &threadName)

    if err != nil {
	glog.Error("INSERT FAILED")
        return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `017`, true)
    }

    a := struct{
		    Id int
		    Name string
		}{
		    threadId,
		    threadName,
		}
    bytes, err1 := json.Marshal(api_request{`ok`, nil, a })
    if err1 != nil {
        return []byte{}, xerrors.NewUIErr(err1.Error(), err1.Error(), `018`, true)
    }

    return bytes, nil
}




