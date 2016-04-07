// Contains *all* custom structs
package main

import "time"

type app struct {
}

type imageBoardClusters struct {
	ID              int
	Descr           string
	LongDescr       string
	BoardLimitCount int
}

type boards struct {
	ID                        int
	Name                      string
	Descr                     string
	ImageBoardClusterID       string
	MaxThreadCount            int  //to be checked in insert thread
	MaxActiveThreadCount      int  //to be checked  in insert thread
	MaxPostsPerThread         int  // to be checked in insert thread
	AreAttachmentsAllowed     bool // to be checked in insert post
	PostLimitsReachedActionID int  // to be checked in insert post
}

type threads struct {
	ID                      int
	Name                    string
	Descr                   string
	boardID                 int
	MaxPostsPerThread       int
	AreAttachmentsAllowed   bool
	LimitsReachedActionID   int
	PostCount               int
	PostCountWithAttachment int
	MinPostLength           int // to be checked in insert post
	MaxPostLength           int // to be checked in insert post
}

type threadPosts struct {
	ID            int
	Body          string
	threadID      int
	attachmentURL *string
	InsertedAt    time.Time
	SourceIP      *string
}

type threadLimitsReachedActions struct {
	ID    int
	Name  string
	Descr string
}

type apiRequest struct {
	Status  string
	Msg     *string
	Payload interface{}
}
