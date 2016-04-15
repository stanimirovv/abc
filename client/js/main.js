    if (! "onhashchange" in window) {
        alert("Please upgrade your browser!");
    }

    function locationHashChanged() {
        //console.log("Change: " + location.hash);
        console.log(arguments.callee.name);
        var path = location.hash.split('/');
        console.log("Path: ", path);
        if(path.length === 1) {
            console.log("locationHashChanged: getBoards()");
            getBoards();
        } else if(path.length === 2) {
            var board = path[1].split(':');

            if( board.length !== 2){
                uiError('Bad path!');
            } else {
                console.log("locationHashChanged: getActiveThreadsForBoard()");
                getActiveThreadsForBoard(board[1]);
            }
        } else if (path.length === 3) {
            var board = path[1].split(':');

            if( board.length !== 2){
                uiError('Bad path!');
                return
            }

            var thread = path[2].split(':');
            if( thread.length !== 2){
                uiError('Bad path!');
            } else {
                console.log("locationHashChanged: getPostsForThread()");
                getPostsForThread(board[1] ,thread[1]);
            }
        } else {
            console.log("Unknown path!");
        }
    }
    window.onhashchange = locationHashChanged;
    window.onload = locationHashChanged;

 boards = {}
 threads = {}
 function getTestPromise(){
    new Promise(function(resolve, reject) {
    $.ajax({
              url: "http://127.0.0.1:8089/api?command=getBoards&api_key=d3c3f756aff00db5cb063765b828e87b",
              type: "GET",
              success: function(){resolve({status: "ok"});},
              error: function(){reject("ERROR!");}
          });

})
.then(function(e) { console.log('done', e); })
.catch(function(e) { console.log('catch: ', e); });
}

function getBoards(resolve){
    console.log("Inside getBoards");
        $.ajax({
                url: "http://127.0.0.1:8089/api?command=getBoards&api_key=d3c3f756aff00db5cb063765b828e87b",
                type: "GET",
                success: function(resp){
                    boards = JSON.parse(resp);
                    if( boards.Status !== 'ok') {
                        uiError(resp.Msg);
                    }

                    var html = '';
                    for (var i = 0; i < boards.Payload.length; i++){
                        console.log(boards.Payload[i]);
                        html += '<h2>'+ boards.Payload[i].Name +'</h2>';
                    }
                if(resolve !== undefined){
                    resolve();
                }

                $("#app").html(html);
              },
              error: function(){}
          });
}


function uiError(errorText){
    alert(errorText);
}

function getActiveThreadsForBoard(boardId){
    console.log("boards: ", boards);
    console.log(boardId);
    if(boards.Status === undefined){
        console.log('MUST LOAD CACHE');
        new Promise(function(resolve, reject) {
            getBoards()
            resolve('ok');
        })
        .then(function(e) { getActiveThreadsForBoardA(boardId) },
                function(e) { console.log('catch: ', e); });
    } else {
        getActiveThreadsForBoardA(boardId);
    }

}

function getActiveThreadsForBoardA(boardId, resolve){
    console.log(arguments.callee.name);
    $.ajax({
          url: "http://127.0.0.1:8089/api?command=getActiveThreadsForBoard&api_key=d3c3f756aff00db5cb063765b828e87b&board_id=" + boardId,
          type: "GET",
          success: function(resp){
              threads = JSON.parse(resp);
              if( threads.Status !== 'ok') {
                  alert(resp.Msg);
              }
              var html = '';
              for(var i = 0; i < boards.Payload.length; i++){
                  if(boards.Payload[i].ID == boardId){
                      html += '<h1>' + boards.Payload[i].Name +'</h1>';
                  }
              }
              for (var i = 0; i < threads.Payload.length; i++){
                  html += '<h2>'+ threads.Payload[i].Name +'</h2>';
              }
              $("#app").html(html);
              if(resolve !== undefined){
                  console.log('resolving the promise');
                  resolve('ok');
              }
              console.log('END getActiveThreadsForBoardA');
          },
          error: function(){}
      });
}

function getPostsForThread(boardId, threadId){
    console.log(arguments.callee.name);
    if(boards.Status === undefined){
        console.log('MUST LOAD CACHE');

        new Promise(function(resolve, reject) {
            getBoards(resolve);})
            .then(function(e) { return new Promise(function(resolve) {getActiveThreadsForBoardA(boardId, resolve);})}, function(e) { console.log('catch: ', e); })
        .then(function(e){getPostsForThreadA(boardId, threadId);}, function(e) {console.log('catch: ', e); })
    } else {
        getPostsForThreadA(boardId, threadId);
    }
}


function getPostsForThreadA(boardId, threadId){
    console.log(arguments.callee.name);

    $.ajax({
              url: "http://127.0.0.1:8089/api?command=getPostsForThread&api_key=d3c3f756aff00db5cb063765b828e87b&thread_id=" + threadId,
              type: "GET",
              success: function(resp ){
                        respObj = JSON.parse(resp);
                        if( respObj.Status !== 'ok') {
                            alert(resp.Msg);
                        }

                        var html = '';
                        console.log("inside getPostsForThreadA: threads: ", threads);
                        for (var i = 0; i < threads.Payload.length; i++){
                            if(threadId == threads.Payload[i].ID){
                                html += '<h2>'+ threads.Payload[i].Name +'</h2>';
                            }
                        }

                        if(respObj.Payload !== null && respObj.Payload.length !== null ){
                            for(var i =0; i < respObj.Payload.length; i++){
                                html += '<div class="postBox">' + respObj.Payload[i].Body + '</div>';
                            }
                        }
                        $("#app").html(html);
                    },
              error: function(){reject("ERROR!");}
          });
}


//xmlhttp.open("GET", "http://localhost:8089/api?command=addPostToThread&api_key=d3c3f756aff00db5cb063765b828e87b&thread_id=" + threadId +
                        //"&thread_post_body=" + escape(document.getElementById('newPostTextArea').value) + "&attachment_url="+escape(attachmentUrl));
//xmlhttp.open("GET", "http://localhost:8089/api?command=addThread&api_key=d3c3f756aff00db5cb063765b828e87b&board_id=" +boardId +"&thread_name="+ escape(document.getElementById('newThreadTextArea').value) );
//xmlhttp1.open("GET", "http://localhost:8089/api?command=addPostToThread&api_key=d3c3f756aff00db5cb063765b828e87b&thread_id=" + addThreadResp.Payload.Id +
