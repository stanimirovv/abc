    if (! "onhashchange" in window) {
        alert("Please upgrade your browser!");
    }

    function locationHashChanged() {
        //console.log("Change: " + location.hash);
        path = location.hash.split('/');
        console.log("Path: ", path);
        if(path.length === 1) {
            getBoards();
        } else if(path.length === 2) {
            board = path[1].split(':');

            if( board.length !== 2){
                uiError('Bad path!');
            } else {
                getActiveThreadsForBoard(board[1]);
            }
        } else if (path.length === 3) {
            board = path[1].split(':');

            if( board.length !== 2){
                uiError('Bad path!');
                return
            }

            thread = path[1].split[':'];
            if( params.length !== 2){
                uiError('Bad path!');
            } else {
                getPostsForThread(board[1] ,thread[1]);
            }
        } else {
            console.log("Unknown path!");
        }
    }
    window.onhashchange = locationHashChanged;
    window.onload = locationHashChanged;

 boards = {}
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

function getBoards(){
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
        .then(function(e) { getActiveThreadsForBoardA(boardId) })
        .catch(function(e) { console.log('catch: ', e); });
    }

}

function getActiveThreadsForBoardA(boardId){

    $.ajax({
          url: "http://127.0.0.1:8089/api?command=getActiveThreadsForBoard&api_key=d3c3f756aff00db5cb063765b828e87b&board_id=" + boardId,
          type: "GET",
          success: function(resp){
              respObj = JSON.parse(resp);
              if( respObj.Status !== 'ok') {
                  alert(resp.Msg);
              }

              var html = '';
              for(var i = 0; i < boards.Payload.length; i++){
                  if(boards.Payload[i].ID == boardId){
                      html += '<h1>' + boards.Payload[i].Name +'</h1>';
                  }
              }
              for (var i = 0; i < respObj.Payload.length; i++){
                  console.log(respObj.Payload[i]);
                  html += '<h2>'+ respObj.Payload[i].Name +'</h2>';
              }
              $("#app").html(html);
          },
          error: function(){}
      });
}

function getPostsForThread(boardId, threadId){
    console.log(boardId);
    if(boards.Status === undefined){
        console.log('MUST LOAD CACHE');
        new Promise(function(resolve, reject) {
            getBoards()
            resolve('ok');
        })
        .then(function(e) { getActiveThreadsForBoardA(boardId) })
        .catch(function(e) { console.log('catch: ', e); });
    }
}

function getPostsForThreadA(boardId, threadId){


}

//xmlhttp.open("GET", "http://127.0.0.1:8089/api?command=getPostsForThread&api_key=d3c3f756aff00db5cb063765b828e87b&thread_id=" + threadId);
//xmlhttp.open("GET", "http://localhost:8089/api?command=addPostToThread&api_key=d3c3f756aff00db5cb063765b828e87b&thread_id=" + threadId +
                        //"&thread_post_body=" + escape(document.getElementById('newPostTextArea').value) + "&attachment_url="+escape(attachmentUrl));

//xmlhttp1.open("GET", "http://127.0.0.1:8089/api?command=getActiveThreadsForBoard&api_key=d3c3f756aff00db5cb063765b828e87b&board_id=" + boardId);
    //xmlhttp1.open("GET", "http://127.0.0.1:8089/api?command=getActiveThreadsForBoard&api_key=d3c3f756aff00db5cb063765b828e87b&board_id=" + boardId);
//xmlhttp.open("GET", "http://localhost:8089/api?command=addThread&api_key=d3c3f756aff00db5cb063765b828e87b&board_id=" +boardId +"&thread_name="+ escape(document.getElementById('newThreadTextArea').value) );
//xmlhttp1.open("GET", "http://localhost:8089/api?command=addPostToThread&api_key=d3c3f756aff00db5cb063765b828e87b&thread_id=" + addThreadResp.Payload.Id +
        //"&thread_post_body=" + escape(document.getElementById('newThreadPostTextArea').value) + "&attachment_url="+escape(document.getElementById('newPostAttachUrlInp').value));
