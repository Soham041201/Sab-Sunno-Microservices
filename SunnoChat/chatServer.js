
userSocketMapping= {}

function chatServer(socket,io){
    socket.on('chat-join',(data)=>{
        console.log(data + socket.id);
        userSocketMapping[socket.id] = socket.id;
        console.log(users);
        users.map((user)=>{
            if(user != socket.id){
                io.to(user).emit('chat-connected',{
                    connectedUser:socket.id,
                    message: "connected"
                });
            }
        })
    });
    socket.on("disconnect", () => {
      console.log("============Socket disconnected=============");
    });
}

module.exports = chatServer;