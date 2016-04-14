    if (! "onhashchange" in window) {
        alert("Please upgrade your browser!");
    }

    function locationHashChanged() {
        //console.log("Change: " + location.hash);
        path = location.hash.split('/');
        console.log("Path: ", path);
        if(path.length === 1) {
            app.getBoards();
        } else if(path.length === 2) {
            board = path[1].split(':');

            if( board.length !== 2){
                app.uiError('Bad path!');
            } else {
                app.getActiveThreadsForBoard(board[1]);
            }
        } else if (path.length === 3) {
            board = path[1].split(':');

            if( board.length !== 2){
                app.uiError('Bad path!');
                return
            } else {
                app.getActiveThreadsForBoard(board[1]);
            }

            thread = path[1].split[':'];
            if( params.length !== 2){
                app.uiError('Bad path!');
            } else {
                app.getPostsForThread(board[1] ,thread[1]);
            }
        } else {
            console.log("Unknown path!");
        }
/*
        if (location.hash === "#boards") {
            obj.getBoards();
        }
        else if (location.hash.indexOf("#thread:") > -1 ){
            obj.getPostsForThread(location.hash.split(":")[2]);
        }
        else if (location.hash.indexOf("#board:") > -1 ){
            obj.showThreadsForBoard(location.hash.split(":")[1]);
        }
        else if (location.hash.indexOf("#new_thread:") > -1 ){
            obj.loadNewThreadTemplate(location.hash.split(":")[1]);
        }
        else{
            document.getElementById("app").innerHTML = "404 page not found";
        }
*/
    }
    window.onhashchange = locationHashChanged;
var app = {

getTestPromise : function(){
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
},

getBoards : function(){
        $.ajax({    url: "http://127.0.0.1:8089/api?command=getBoards&api_key=d3c3f756aff00db5cb063765b828e87b",
              type: "GET",
              success: function(){},
              error: function(){}
          });
},


uiError : function(errorText){
    alert(errorText);
},

getActiveThreadsForBoard : function(boardId){
    console.log(boardId);
    $.ajax({
          url: "http://127.0.0.1:8089/api?command=getActiveThreadsForBoard&api_key=d3c3f756aff00db5cb063765b828e87b&board_id=" + boardId,
          type: "GET",
          success: function(){},
          error: function(){}
      });
}

}
//xmlhttp.open("GET", "http://127.0.0.1:8089/api?command=getBoards&api_key=d3c3f756aff00db5cb063765b828e87b");
//xmlhttp.open("GET", "http://127.0.0.1:8089/api?command=getPostsForThread&api_key=d3c3f756aff00db5cb063765b828e87b&thread_id=" + threadId);
//xmlhttp.open("GET", "http://localhost:8089/api?command=addPostToThread&api_key=d3c3f756aff00db5cb063765b828e87b&thread_id=" + threadId +
                        //"&thread_post_body=" + escape(document.getElementById('newPostTextArea').value) + "&attachment_url="+escape(attachmentUrl));

//xmlhttp1.open("GET", "http://127.0.0.1:8089/api?command=getActiveThreadsForBoard&api_key=d3c3f756aff00db5cb063765b828e87b&board_id=" + boardId);
    //xmlhttp1.open("GET", "http://127.0.0.1:8089/api?command=getActiveThreadsForBoard&api_key=d3c3f756aff00db5cb063765b828e87b&board_id=" + boardId);
//xmlhttp.open("GET", "http://localhost:8089/api?command=addThread&api_key=d3c3f756aff00db5cb063765b828e87b&board_id=" +boardId +"&thread_name="+ escape(document.getElementById('newThreadTextArea').value) );
//xmlhttp1.open("GET", "http://localhost:8089/api?command=addPostToThread&api_key=d3c3f756aff00db5cb063765b828e87b&thread_id=" + addThreadResp.Payload.Id +
        //"&thread_post_body=" + escape(document.getElementById('newThreadPostTextArea').value) + "&attachment_url="+escape(document.getElementById('newPostAttachUrlInp').value));
