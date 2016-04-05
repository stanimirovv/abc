package main

import (
	//DB
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/iambc/xerrors"
	"github.com/golang/glog"
)

var dbh *sql.DB
var dbConnString string

//Relational database implementation for writer interface
type writerrdb struct {
    dbh *sql.DB
}

func (db *writerrdb) getBoards(apiKey string) (currBoards []boards, err error) {
    glog.Info(" apiKey: ", apiKey)
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

func (db *writerrdb) getActiveThreadsForBoard(apiKey string, boardId int) (activeThreads []threads, err error) {
    rows, err := dbh.Query(`select t.id, t.name, count(*), (select count(*) from thread_posts where thread_id = t.id and attachment_url is not null) from threads t  
				join boards b on b.id = t.board_id 
				join image_board_clusters ibc on ibc.id = b.image_board_cluster_id 
				left join thread_posts tp on tp.thread_id = t.id
			    where t.is_active = TRUE and t.board_id = $1 and ibc.api_key = $2 group by 1,2 order by t.id;`, boardId, apiKey)
    if err != nil {
        return activeThreads, xerrors.NewUIErr(err.Error(), err.Error(), `006`, true)
    }
    defer rows.Close()

    for rows.Next() {
	glog.Info("Popped new thread")
        var thread threads
        err = rows.Scan(&thread.Id, &thread.Name, &thread.PostCount, &thread.PostCountWithAttachment)
        if err != nil {
            return activeThreads, xerrors.NewUIErr(err.Error(), err.Error(), `007`, true)
        }
        activeThreads = append(activeThreads, thread)
    }
    return activeThreads, nil
}

func (db *writerrdb) getPostsForThread(apiKey string, threadId int) (currPosts []thread_posts, err error) {
    rows, err := dbh.Query(`select tp.id, tp.body, tp.attachment_url, tp.inserted_at, tp.source_ip 
			    from thread_posts tp join threads t on t.id = tp.thread_id 
						 join boards b on b.id = t.board_id 
						 join image_board_clusters ibc on ibc.id = b.image_board_cluster_id 
			    where tp.thread_id = $1 and ibc.api_key = $2 and t.is_active = true;`, threadId, apiKey)
    if err != nil {
	glog.Error(err)
        return currPosts, xerrors.NewSysErr()
    }
    defer rows.Close()

    for rows.Next() {
	glog.Info("new post for thread with id: ", threadId)
        var currPost thread_posts
        err = rows.Scan(&currPost.Id, &currPost.Body, &currPost.AttachmentUrl, &currPost.InsertedAt, &currPost.SourceIp)
        if err != nil {
	    glog.Error(err)
            return currPosts, xerrors.NewSysErr()
        }
        currPosts = append(currPosts, currPost)
    }
    return currPosts, err
}

func (db *writerrdb) addPostToThread(threadId int, threadBodyPost string, attachmentUrl *string, clientRemoteAddr string) (err error) {
    _, err = dbh.Query("INSERT INTO thread_posts(body, thread_id, attachment_url, source_ip) VALUES($1, $2, $3, $4)", threadBodyPost, threadId, attachmentUrl, clientRemoteAddr)

    if err != nil {
	glog.Error(err)
        return xerrors.NewUIErr(err.Error(), err.Error(), `011`, true)
    }
    return nil
}

func (db *writerrdb) addThread(boardId int, threadName string) (threads, error) {

    var threadId int
    err := dbh.QueryRow("INSERT INTO threads(name, board_id, limits_reached_action_id, max_posts_per_thread) VALUES($1, $2, 1, 10)  RETURNING id, name", threadName, boardId).Scan(&threadId, &threadName)

    if err != nil {
	glog.Error("INSERT FAILED")
        return threads{Id:-1, Name:`err`}, xerrors.NewUIErr(err.Error(), err.Error(), `017`, true)
    }
    return threads{Id:threadId, Name:threadName}, nil
}

func (db *writerrdb) isThreadLimitReached(boardId int) (bool, error) {
    var isLimitReached bool
    err := dbh.QueryRow("select (select count(*) from threads  where board_id = $1) > thread_setting_max_thread_count  from boards where id = $1;", boardId).Scan(&isLimitReached)
    if err != nil {
	glog.Error("COULD NOT SELECT thread_count")
	return true, xerrors.NewUIErr(err.Error(), err.Error(), `015`, true)
    }

    return isLimitReached, nil
}


func (db *writerrdb) isPostLimitReached(threadId int) (bool, threads, error){
    var isLimitReached bool
    var thread threads
    err := dbh.QueryRow("select (select count(*) from thread_posts  where thread_id = $1) > max_posts_per_thread, min_post_length, max_post_length  from threads where id = $1;", threadId).Scan(&isLimitReached, &thread.MinPostLength, &thread.MaxPostLength)
    if err != nil {
	return true, thread, xerrors.NewUIErr(err.Error(), err.Error(), `009`, true)
    }
    return isLimitReached, thread, err
}
