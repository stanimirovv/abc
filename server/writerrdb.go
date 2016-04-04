package main

import (
	//DB
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/iambc/xerrors"
)

var dbh *sql.DB
var dbConnString string

//Relational database implementation for writer interface
type writerrdb struct {
}

func (db *writerrdb) getBoards(apiKey string) (currBoards []boards, err error) {
    rows, err := dbh.Query("select b.id, b.name, b.descr from boards b join image_board_clusters ibc on ibc.id = b.image_board_cluster_id where api_key = $1;", apiKey)
    if err != nil {
	return currBoards, xerrors.NewUIErr(err.Error(), err.Error(), `002`, true)
    }
    defer rows.Close()

    for rows.Next() {
	var board boards
	err = rows.Scan(&board.Id, &board.Name, &board.Descr)
	if err != nil {
	    return currBoards, xerrors.NewUIErr(err.Error(), err.Error(), `003`, true)
	}
	currBoards = append(currBoards, board)
    }
    return currBoards, nil
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
