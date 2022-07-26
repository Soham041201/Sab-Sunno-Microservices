
function chatServer(socket){
    socket.on('chat-join',(data)=>{
        console.log(data);
    });
}

module.exports = chatServer;