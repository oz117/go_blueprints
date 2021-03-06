$(function(){
  var socket = null;
  var msgBox = $("#chatbox textarea");
  var messages = $("#messages");

  $("#chatbox").submit(function(){
    if (!msgBox.val()) return false;
    if (!socket) {
      alert("Error: There is no socket connection.");
      return false;
    }
    socket.send(JSON.stringify({"Message": msgBox.val()}));
    msgBox.val("");
    return false;
  });
  if (!window["WebSocket"]) {
    alert("Error: Your browser does not support web sockets.");
  } else {
    socket = new WebSocket("ws://"+document.location.host+"/room");
    socket.onclose = function() {
      alert("Connection has been closed.");
    }
    socket.onmessage = function(e) {
      var msg = eval("(" + e.data + ")");
      messages.append(
        $("<li>").append(
          $("<strong>").text(msg.Name + ": "),
          $("<span>").text(msg.Message)
        )
      );
    }
  }
});
