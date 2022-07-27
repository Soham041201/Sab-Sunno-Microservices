
users= []

function chatServer(socket,io){
    socket.on('chat-join',(data)=>{
        console.log(data + socket.id);
        !users.includes(socket.id) && users.push(socket.id);
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