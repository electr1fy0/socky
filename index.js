const WS_URL = "ws://localhost:8080/ws";

const socket = new WebSocket(WS_URL);

socket.onopen = () => {
  console.log("connected to server");
};

socket.onmessage = (event) => {
  console.log(event.data);
  let x = JSON.parse(event.data);
  console.log(x.gridString);
};

socket.onerror = (event) => {
  console.log(event);
};

document.addEventListener("keydown", (event) => {
  console.log(event);
  socket.send(event.key);
});
