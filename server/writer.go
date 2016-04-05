package main

type Writer interface {
    getBoards(apiKey string) (currBoards []boards, err error)
    getActiveThreadsForBoard(apiKey string, boardId int) (activeThreads []threads, err error)
    getPostsForThread(apiKey string, threadId int) (currPosts []thread_posts, err error)
    addPostToThread(threadId int, threadBodyPost string, attachmentUrl *string, clientRemoteAddr string) (err error)
    addThread(boardId int, threadName string) (threads, error)
    isThreadLimitReached(boardId int) (bool, error)
    isPostLimitReached(threadId int) (bool, threads, error)
}
