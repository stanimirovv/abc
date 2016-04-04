// Contains *all* custom structs
package main

import "time"

type app struct {

}

type image_board_clusters struct {
    Id int
    Descr string
    LongDescr string
    BoardLimitCount int
}

type boards struct {
    Id int
    Name string
    Descr string
    ImageBoardClusterId string
    MaxThreadCount int //to be checked in insert thread
    MaxActiveThreadCount int //to be checked  in insert thread
    MaxPostsPerThread int // to be checked in insert thread
    AreAttachmentsAllowed bool // to be checked in insert post
    PostLimitsReachedActionId int // to be checked in insert post
}

type threads struct{
    Id int
    Name string
    Descr string
    BoardId int
    MaxPostsPerThread int
    AreAttachmentsAllowed bool
    LimitsReachedActionId int
    PostCount int
    PostCountWithAttachment int
    MinPostLength int // to be checked in insert post
    MaxPostLength int // to be checked in insert post
}

type thread_posts struct{
    Id int
    Body string
    ThreadId int
    AttachmentUrl *string
    InsertedAt time.Time
    SourceIp *string
}

type thread_limits_reached_actions struct{
    Id	    int
    Name    string
    Descr   string
}

type api_request struct{
    Status  string
    Msg	    *string
    Payload interface{}
}
