
users= new Set()

function chatServer(socket,io,users){
    socket.on('chat-join',(data)=>{
        console.log(data + socket.id);
        users.add(socket.id)
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
}

module.exports = chatServer;