package main

import (
	//DB
	"database/sql"

	"github.com/golang/glog"
	"github.com/iambc/xerrors"
	_ "github.com/lib/pq"
)

//Relational database implementation for writer interface
type writerrdb struct {
	*sql.DB
}

func (db *writerrdb) getBoards(apiKey string) (currBoards []boards, err error) {
	glog.Info(" apiKey: ", apiKey)
	rows, err := db.Query("select b.ID, b.name, b.descr from boards b join image_board_clusters ibc on ibc.ID = b.image_board_cluster_ID where api_key = $1;", apiKey)
	if err != nil {
		return currBoards, xerrors.NewUIErr(err.Error(), err.Error(), `002`, true)
	}
	defer rows.Close()

	for rows.Next() {
		var board boards
		err = rows.Scan(&board.ID, &board.Name, &board.Descr)
		if err != nil {
			return currBoards, xerrors.NewUIErr(err.Error(), err.Error(), `003`, true)
		}
		currBoards = append(currBoards, board)
	}
	return currBoards, nil
}

func (db *writerrdb) getActiveThreadsForBoard(apiKey string, boardID int) (activeThreads []threads, err error) {
	rows, err := db.Query(`select t.ID, t.name, count(*), (select count(*) from thread_posts where thread_ID = t.ID and attachment_url is not null) from threads t
				join boards b on b.ID = t.board_ID
				join image_board_clusters ibc on ibc.ID = b.image_board_cluster_ID
				left join thread_posts tp on tp.thread_ID = t.ID
			    where t.is_active = TRUE and t.board_ID = $1 and ibc.api_key = $2 group by 1,2 order by t.ID;`, boardID, apiKey)
	if err != nil {
		return activeThreads, xerrors.NewUIErr(err.Error(), err.Error(), `006`, true)
	}
	defer rows.Close()

	for rows.Next() {
		var thread threads
		err = rows.Scan(&thread.ID, &thread.Name, &thread.PostCount, &thread.PostCountWithAttachment)
		if err != nil {
			return activeThreads, xerrors.NewUIErr(err.Error(), err.Error(), `007`, true)
		}
		activeThreads = append(activeThreads, thread)
	}
	return activeThreads, nil
}

func (db *writerrdb) getPostsForThread(apiKey string, threadID int) (currPosts []threadPosts, err error) {
	rows, err := db.Query(`select tp.ID, tp.body, tp.attachment_url, tp.inserted_at, tp.source_ip
			    from thread_posts tp join threads t on t.ID = tp.thread_ID
						 join boards b on b.ID = t.board_ID
						 join image_board_clusters ibc on ibc.ID = b.image_board_cluster_ID
			    where tp.thread_ID = $1 and ibc.api_key = $2 and t.is_active = true;`, threadID, apiKey)
	if err != nil {
		glog.Error(err)
		return currPosts, xerrors.NewSysErr()
	}
	defer rows.Close()

	for rows.Next() {
		glog.Info("new post for thread with ID: ", threadID)
		var currPost threadPosts
		err = rows.Scan(&currPost.ID, &currPost.Body, &currPost.attachmentURL, &currPost.InsertedAt, &currPost.SourceIP)
		if err != nil {
			glog.Error(err)
			return currPosts, xerrors.NewSysErr()
		}
		currPosts = append(currPosts, currPost)
	}
	return currPosts, err
}

func (db *writerrdb) addPostToThread(threadID int, threadBodyPost string, attachmentURL *string, clientRemoteAddr string) (err error) {
	_, err = db.Query("INSERT INTO thread_posts(body, thread_ID, attachment_url, source_ip) VALUES($1, $2, $3, $4)", threadBodyPost, threadID, attachmentURL, clientRemoteAddr)

	if err != nil {
		glog.Error(err)
		return xerrors.NewUIErr(err.Error(), err.Error(), `011`, true)
	}
	return nil
}

func (db *writerrdb) addThread(boardID int, threadName string) (threads, error) {

	var threadID int
	err := db.QueryRow("INSERT INTO threads(name, board_ID, limits_reached_action_ID, max_posts_per_thread) VALUES($1, $2, 1, 10)  RETURNING ID, name", threadName, boardID).Scan(&threadID, &threadName)

	if err != nil {
		glog.Error("INSERT FAILED")
		return threads{ID: -1, Name: `err`}, xerrors.NewUIErr(err.Error(), err.Error(), `017`, true)
	}
	return threads{ID: threadID, Name: threadName}, nil
}

func (db *writerrdb) isThreadLimitReached(boardID int) (bool, error) {
	var isLimitReached bool
	err := db.QueryRow("select (select count(*) from threads  where board_ID = $1) > thread_setting_max_thread_count  from boards where ID = $1;", boardID).Scan(&isLimitReached)
	if err != nil {
		glog.Error("COULD NOT SELECT thread_count")
		return true, xerrors.NewUIErr(err.Error(), err.Error(), `015`, true)
	}

	return isLimitReached, nil
}

func (db *writerrdb) isPostLimitReached(threadID int) (bool, threads, error) {
	var isLimitReached bool
	var thread threads
	err := db.QueryRow("select (select count(*) from thread_posts  where thread_ID = $1) > max_posts_per_thread, min_post_length, max_post_length  from threads where ID = $1;", threadID).Scan(&isLimitReached, &thread.MinPostLength, &thread.MaxPostLength)
	if err != nil {
		return true, thread, xerrors.NewUIErr(err.Error(), err.Error(), `009`, true)
	}
	return isLimitReached, thread, err
}

func (db *writerrdb) archiveThread(threadID int) {
	db.QueryRow("UPDATE threads set is_active = false where ID = $1", threadID).Scan()
}
