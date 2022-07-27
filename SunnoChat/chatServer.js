const users = []

function chatServer(socket,io){
    socket.on('chat-join',(data)=>{
        console.log(data + socket.id);
        users.push(socket.id);
        console.log(users);
    });
}

module.exports = chatServer;