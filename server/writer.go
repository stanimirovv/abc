package main

type Writer interface {
    getBoards(apiKey string) []boards
    getActiveThreadsForBoard(apiKey string, boardId int) []threads
    getPostsForThread(apiKey string, threadId int) []thread_posts
    addPostToThread(threadId int, threadBodyPost string, attachmentUrl *string, clientRemoteAddr string) error
    addThread(boardId int, threadName string) (threads, error)
    isThreadLimitReached(boardId int) (int, error)
    isPostLimitReached(threadId int) (int, error)
}
