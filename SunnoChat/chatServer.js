function chatServer(socket, io) {
  const clients = Array.from(io.sockets.adapter.rooms.get("chat-room") || []);

  socket.on("chat-join", (data) => {
    console.log(data + socket.id);
    socket.join("chat-room");
    console.log("Clients connected", clients);
    clients.map((user) => {
      if (user != socket.id) {
        io.to(user).emit("chat-connected", {
          connectedUser: socket.id,
          message: "connected",
        });
      }
    });
  });

  socket.on("send-message", (data) => {
    console.log(data);
    io.to(socket.id).emit("receive-message", {
      message: data,
      sender: socket.id,
    });
  })


  socket.on("disconnect", () => {
    console.log("============Socket disconnected=============", socket.id);
    socket.leave("chat-room");
  });
}

module.exports = chatServer;
