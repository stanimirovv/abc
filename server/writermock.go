package main

import (
	"github.com/iambc/xerrors"
)

type writermock struct {
}

func (db *writermock) getBoards(apiKey string) (currBoards []boards, err error) {
	if apiKey == `err` {
		return currBoards, xerrors.NewUIErr(`test err`, `test err`, `002`, true)
	}

	currBoards = append(currBoards, boards{Id: 1, Name: "Mock Board"})
	return currBoards, nil
}

func (db *writermock) getActiveThreadsForBoard(apiKey string, boardID int) (activeThreads []threads, err error) {
	if apiKey == `err` {
		return activeThreads, xerrors.NewUIErr(`test err`, `test err`, `002`, true)
	} else if apiKey == `empty` {
		return activeThreads, nil
	}
	activeThreads = append(activeThreads, threads{Id: 1, Name: `Mock Thread`})
	return activeThreads, nil

}

func (db *writermock) getPostsForThread(apiKey string, threadID int) (currPosts []thread_posts, err error) {
	if apiKey == `err` {
		return currPosts, xerrors.NewUIErr(`test err`, `test err`, `002`, true)
	} else if apiKey == `empty` {
		return currPosts, nil
	}
	currPosts = append(currPosts, thread_posts{})
	return currPosts, nil
}

func (db *writermock) addPostToThread(threadID int, threadBodyPost string, attachmentUrl *string, clientRemoteAddr string) (err error) {
	if threadID == 1 {
		return xerrors.NewUIErr(`test err`, `test err`, `002`, true)
	}
	return nil
}

func (db *writermock) addThread(boardID int, threadName string) (threads, error) {
	if boardID == 1 {
		return threads{}, xerrors.NewUIErr(`test err`, `test err`, `002`, true)
	}
	return threads{Id: 1, boardID: boardID, Name: threadName}, nil
}

func (db *writermock) isThreadLimitReached(boardID int) (bool, error) {
	if boardID == 2 {
		return true, nil
	} else if boardID == 3 {
		return false, nil
	} else {
		return false, xerrors.NewUIErr(`test err`, `test err`, `002`, true)
	}
}

func (db *writermock) isPostLimitReached(threadID int) (bool, threads, error) {
	if threadID == 2 {
		return true, threads{Id: threadID, Name: `Mock Thread`}, nil
	} else if threadID == 3 {
		return false, threads{Id: threadID, Name: `Mock Thread`, MinPostLength: 2, MaxPostLength: 4}, nil
	} else {
		return false, threads{Id: threadID, Name: `Mock Thread`}, xerrors.NewUIErr(`test err`, `test err`, `002`, true)
	}
}

func (db *writermock) archiveThread(threadID int) {}
