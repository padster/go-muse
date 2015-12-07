IMG = null;
zoom = 100;

randRange = function(lb, ub) {
  return lb + Math.random() * (ub - lb);
};

handleBlink = function() {
  var growth = randRange(1.3, 1.9);

  // Using zoom:
  // zoom = (zoom * growth);
  // document.body.style.zoom = zoom + '%';
  
  // Using width:
  IMG.style.width = (IMG.width * growth) + 'px';
};

openSocket = function() {
  var socket = new WebSocket("ws://localhost:8888/sock");
  socket.onopen = function(e) { console.log("Connected."); };
  socket.onmessage = function(e) {
    switch (e.data | 0) {
      case 0:
        handleBlink();
        break;
    }
  };
};

$(document).ready(function() {
  IMG = document.getElementById('img');

  console.log("Opening websocket channel...");
  openSocket();
});