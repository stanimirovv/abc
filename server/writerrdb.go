package main

//Relational database implementation for writer interface
type writerrdb struct {
}

func (db *writerrdb) getBoards(apiKey string) []boards {
}

func (db *writerrdb) getActiveThreadsForBoard(apiKey string, boardId int) []threads {

}

func (db *writerrdb) getPostsForThread(apiKey string, threadId int) []thread_posts {

}

func (db *writerrdb) addPostToThread(threadId int, threadBodyPost string, attachmentUrl *string, clientRemoteAddr string) error {
}

func (db *writerrdb) addThread(boardId int, threadName string) (threads, error) {

}

func (db *writerrdb) getThreadCount(boardId int) (int, error) {

}
