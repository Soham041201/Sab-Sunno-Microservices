function chatServer(socket, io) {
  const clients = Array.from(io.sockets.adapter.rooms.get("chat-room") || []);


  socket.on("chat-join", (data) => {
    console.log(data + socket.id);
    socket.join("chat-room");
    console.log("Clients connected", Array.from(io.sockets.adapter.rooms.get("chat-room")) );
    Array.from(io.sockets.adapter.rooms.get("chat-room")).map((user) => {
      if (user != socket.id) {
        io.to(user).emit("chat-connected", JSON.stringify({
          connectedUser: socket.id,
          message: "connected",
        }));
        socket.emit("chat-connected", JSON.stringify({
            connectedUser: user,
            message: "connected",
        }));
      }
    });
  });

  socket.on("send-message", (data) => {
    console.log("=====================send-message====================");
    data = JSON.parse(data);
    console.log({message : data.message, reciever : data.reciever});
    io.to(data.reciever).emit("recieve-message",JSON.stringify({
      message: data.message,
      sender: socket.id,
    }));
  })

  socket.on('is_online',()=>{
    console.log("===============is_online SOCKET EVENT==========");
    console.log("=======online socket========",socket.id);

  })

  socket.on('last_seen',()=>{})



  socket.on('leave-chat',()=>{
    console.log("============Socket left the chat=============",socket.id);
    console.log("Clients",clients);
      socket.leave("chat-room");
  })

  socket.on("disconnect", () => {
    console.log("============Socket disconnected=============", socket.id);
    socket.leave("chat-room");
  });
}

module.exports = chatServer;
