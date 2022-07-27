


function chatServer(socket,io){
    const clients = Array.from(io.sockets.adapter.rooms.get('chat-room') || []);
    
    socket.on('chat-join',(data)=>{
        console.log(data + socket.id);
        socket.join('chat-room');
        console.log(users);
        clients.map((user)=>{
            if(user != socket.id){
                io.to(user).emit('chat-connected',{
                    connectedUser:socket.id,
                    message: "connected"
                });
            }
        })
    });

    socket.on("disconnect", () => {
        console.log("============Socket disconnected=============", socket.id);
        socket.leave('chat-room');
    });
}

module.exports = chatServer;