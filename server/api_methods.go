package main

/*
This file contains all of the API functions.
*/

import (
	"encoding/json"
	"strconv"

	"github.com/golang/glog"
	"github.com/iambc/xerrors"
)

type abcAPI struct {
	storageType string
	wr          Writer
}

func (api *abcAPI) getBoards(apiKey string) ([]byte, error) {

	currBoards, err := api.wr.getBoards(apiKey)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := json.Marshal(apiRequest{"ok", nil, &currBoards})
	if err != nil {
		return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `004`, true)
	}
	return bytes, nil
}

func (api *abcAPI) getActiveThreadsForBoard(apiKey string, boardID int) ([]byte, error) {

	activeThreads, err := api.wr.getActiveThreadsForBoard(apiKey, boardID)
	if err != nil {
		return []byte{}, err
	}
	var bytes []byte
	if len(activeThreads) == 0 {
		errMsg := "No objects returned."
		bytes, err = json.Marshal(apiRequest{"error", &errMsg, &activeThreads})
	} else {
		bytes, err = json.Marshal(apiRequest{"ok", nil, &activeThreads})
	}

	if err != nil {
		return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `008`, true)
	}

	return bytes, nil
}

func (api *abcAPI) getPostsForThread(apiKey string, threadID int) ([]byte, error) {
	currPosts, err := api.wr.getPostsForThread(apiKey, threadID)

	var bytes []byte
	if len(currPosts) == 0 {
		errMsg := "No objects returned."
		var tmp []int
		bytes, err = json.Marshal(apiRequest{"ok", &errMsg, &tmp})
	} else {
		bytes, err = json.Marshal(apiRequest{"ok", nil, &currPosts})
	}

	if err != nil {
		return []byte{}, xerrors.NewSysErr()
	}

	return bytes, nil
}

func (api *abcAPI) addPostToThread(threadID int, threadBodyPost string, attachmentURL *string, clientRemoteAddr string) ([]byte, error) {
	isLimitReached, thr, err := api.wr.isPostLimitReached(threadID)

	if err != nil {
		return []byte{}, err
	}
	if isLimitReached {
		api.wr.archiveThread(threadID)
		return []byte{}, xerrors.NewUIErr(`Thread post limit reached!`, `Thread post limit reached!`, `010`, true)
	}

	if thr.MinPostLength > len(threadBodyPost) && thr.MinPostLength != -1 {
		return []byte{}, xerrors.NewUIErr(`Post length is less than minimum length!`, `Post length is less than minimum length! post length: `+strconv.Itoa(len(threadBodyPost))+` min length: `+strconv.Itoa(thr.MinPostLength), `020`, false)
	}

	if thr.MaxPostLength < len(threadBodyPost) && thr.MaxPostLength != -1 {
		return []byte{}, xerrors.NewUIErr(`Post length is more than maximum length!`, `Post length is more than maximum length! post length: `+strconv.Itoa(len(threadBodyPost))+` max length: `+strconv.Itoa(thr.MaxPostLength), `021`, false)
	}

	err = api.wr.addPostToThread(threadID, threadBodyPost, attachmentURL, clientRemoteAddr)

	if err != nil {
		glog.Error(err)
		return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `011`, true)
	}

	var bytes []byte
	bytes, err = json.Marshal(apiRequest{"ok", nil, nil})
	if err != nil {
		return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `012`, true)
	}

	return bytes, nil
}

func (api *abcAPI) addThread(boardID int, threadName string) ([]byte, error) {

	isLimitReached, err := api.wr.isThreadLimitReached(boardID)

	if err != nil {
		return []byte{}, err
	}

	if isLimitReached {
		return []byte{}, xerrors.NewUIErr(`Thread limit reached!`, `Thread limit reached!`, `016`, true)
	}

	thr, err := api.wr.addThread(boardID, threadName)
	if err != nil {
		glog.Error("INSERT FAILED")
		return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `017`, true)
	}

	a := struct {
		ID   int
		Name string
	}{
		thr.ID,
		thr.Name,
	}
	var bytes []byte
	bytes, err = json.Marshal(apiRequest{`ok`, nil, a})
	if err != nil {
		return []byte{}, xerrors.NewUIErr(err.Error(), err.Error(), `018`, true)
	}

	return bytes, nil
}
