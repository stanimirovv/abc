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
    isLimitReached, thr, err := api.wr.isPostLimitReached(threadId)

    if err != nil {
	return []byte{}, err
    }

    if isLimitReached {
	api.wr.archiveThread(threadId)
	return []byte{}, xerrors.NewUIErr(`Thread post limit reached!`, `Thread post limit reached!`, `010`, true)
    }

   if(thr.MinPostLength > len(threadBodyPost)  && thr.MinPostLength != -1){
	return []byte{}, xerrors.NewUIErr(`Post length is less than minimum length!`, `Post length is less than minimum length! post length: ` + strconv.Itoa(len(threadBodyPost))  +` min length: ` + strconv.Itoa(thr.MinPostLength) , `020`, false)
    }
   if(thr.MaxPostLength < len(threadBodyPost)  && thr.MaxPostLength != -1){
	return []byte{}, xerrors.NewUIErr(`Post length is more than maximum length!`, `Post length is more than maximum length! post length: ` + strconv.Itoa(len(threadBodyPost))  +` max length: ` + strconv.Itoa(thr.MaxPostLength) , `021`, false)
    }

    err = api.wr.addPostToThread(threadId, threadBodyPost, attachmentUrl, clientRemoteAddr)

    if err != nil {
	glog.Error(err)
        return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `011`, true)
    }

    var bytes []byte
    bytes, err = json.Marshal(api_request{"ok", nil, nil})
    if err != nil {
        return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `012`, true)
    }

    return bytes, nil
}


func (api *abc_api) addThread(boardId int, threadName string) ([]byte, error) {

    isLimitReached, err := api.wr.isThreadLimitReached(boardId)

    if err != nil {
	return []byte{}, err
    }

    if isLimitReached {
	return []byte{}, xerrors.NewUIErr(`Thread limit reached!`, `Thread limit reached!`, `016`, true)
    }

    thr, err := api.wr.addThread(boardId, threadName)
    if err != nil {
	glog.Error("INSERT FAILED")
        return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `017`, true)
    }

    a := struct{
		    Id int
		    Name string
		}{
		    thr.Id,
		    thr.Name,
		}
    var bytes []byte
    bytes, err = json.Marshal(api_request{`ok`, nil, a })
    if err != nil {
        return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `018`, true)
    }

    return bytes, nil
}




