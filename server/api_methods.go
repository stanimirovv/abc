package main

/*
This file contains all of the API functions.
*/

import (
"encoding/json"
	"github.com/iambc/xerrors"
	"github.com/golang/glog"
	"strconv"
	)

type abc_api struct {
    storage_type string
    wr Writer
}

func (api *abc_api) getBoards(apiKey string)  ([]byte, error) {

    currBoards, err := api.wr.getBoards(apiKey)
    if err != nil {
	return []byte{}, err
    }
    bytes, err := json.Marshal(api_request{"ok", nil, &currBoards})
    if err != nil {
	return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `004`, true)
    }
    return bytes, nil
}


func (api *abc_api) getActiveThreadsForBoard(apiKey string, boardId int)  ([]byte, error) {

    activeThreads, err := api.wr.getActiveThreadsForBoard(apiKey, boardId)
    if err != nil {
	return []byte{}, err
    }
    var bytes []byte
    if(len(activeThreads) == 0){
        errMsg := "No objects returned."
        bytes, err = json.Marshal(api_request{"error", &errMsg, &activeThreads})
    }else {
        bytes, err = json.Marshal(api_request{"ok", nil, &activeThreads})
    }

    if err != nil {
        return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `008`, true)
    }

    return bytes, nil
}


func (api *abc_api) getPostsForThread(apiKey string, threadId int)  ([]byte, error) {
    currPosts, err := api.wr.getPostsForThread(apiKey, threadId)

    var bytes []byte
    if(len(currPosts) == 0){
	errMsg := "No objects returned."
	bytes, err = json.Marshal(api_request{"error", &errMsg, &currPosts})
    }else {
	bytes, err = json.Marshal(api_request{"ok", nil, &currPosts})
    }

    if err != nil {
        return []byte{}, xerrors.NewSysErr()
    }

    return bytes, nil
}


func (api *abc_api) addPostToThread(threadId int, threadBodyPost string, attachmentUrl *string, clientRemoteAddr string) ([]byte,error) {
    var isLimitReached bool
    var maxPostLength int
    var minPostLength int
    err := dbh.QueryRow("select (select count(*) from thread_posts  where thread_id = $1) > max_posts_per_thread, min_post_length, max_post_length  from threads where id = $1;", threadId).Scan(&isLimitReached, &minPostLength, &maxPostLength)
    if err != nil {
	return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `009`, true)
    }

    if isLimitReached {
	dbh.QueryRow("UPDATE threads set is_active = false where id = $1", threadId).Scan()
	return []byte{}, xerrors.NewUIErr(`Thread post limit reached!`, `Thread post limit reached!`, `010`, true)
    }

   if(minPostLength > len(threadBodyPost)  && minPostLength != -1){
	return []byte{}, xerrors.NewUIErr(`Post length is less than minimum length!`, `Post length is less than minimum length! post length: ` + strconv.Itoa(len(threadBodyPost))  +` min length: ` + strconv.Itoa(minPostLength) , `020`, false)
    }
   if(maxPostLength < len(threadBodyPost)  && maxPostLength != -1){
	return []byte{}, xerrors.NewUIErr(`Post length is more than maximum length!`, `Post length is more than maximum length! post length: ` + strconv.Itoa(len(threadBodyPost))  +` max length: ` + strconv.Itoa(maxPostLength) , `021`, false)
    }

    _, err = dbh.Query("INSERT INTO thread_posts(body, thread_id, attachment_url, source_ip) VALUES($1, $2, $3, $4)", threadBodyPost, threadId, attachmentUrl, clientRemoteAddr)

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


func (api *abc_api) addThread(boardId int, threadName string) ([]byte, error) {
    var isLimitReached bool
    err := dbh.QueryRow("select (select count(*) from threads  where board_id = $1) > thread_setting_max_thread_count  from boards where id = $1;", boardId).Scan(&isLimitReached)
    if err != nil {
	glog.Error("COULD NOT SELECT thread_count")
	return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `015`, true)
    }
    if isLimitReached {
	return []byte{}, xerrors.NewUIErr(`Thread limit reached!`, `Thread limit reached!`, `016`, true)
    }

    var threadId int
    err = dbh.QueryRow("INSERT INTO threads(name, board_id, limits_reached_action_id, max_posts_per_thread) VALUES($1, $2, 1, 10)  RETURNING id, name", threadName, boardId).Scan(&threadId, &threadName)

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




