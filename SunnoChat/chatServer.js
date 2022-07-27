
function chatServer(socket,io){
    socket.on('chat-join',(data)=>{
        console.log(data + socket.id);
    });
}

module.exports = chatServer;