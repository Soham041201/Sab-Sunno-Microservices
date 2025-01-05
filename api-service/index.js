const { db_connection, Room } = require('./db-connect');
const express = require('express');
const app = express();
const cors = require('cors');
const { ObjectId } = require('mongodb');
require('dotenv').config();
const server = require('http').createServer(app);
const chatServer = require('./SunnoChat/chatServer');
const routes = require('./routes');
const WebSocket = require('ws'); // Import the ws library

const {
  fetchEphimerialToken,
  sendLocalDescriptionToOpenAi,
} = require('./openai');

const wss = new WebSocket.Server({ port: 8001, path: '/ws' });
const ioSockets = new Set();

wss.on('connection', (ws) => {
  console.log('Client connected');

  const message = { event: 'connected', data: 'Welcome to sab sunno' };
  ws.send(JSON.stringify(message));

  ws.on('message', (message) => {
    console.log(`Received: ${message}`);

    ioSockets.forEach((socket) => {
      const socketEvent = JSON.parse(message);
      console.log('Socket event from go', socketEvent);

      if (socketEvent.type === 'answer') {
        socket.emit('open-ai-answer', {
          peerId: 'open-ai',
          sessionDescription: socketEvent.sdp,
        }); // Emit a Socket.IO event
      }
      if (socketEvent.event === 'ice-candidate') {
        socket.emit('open-ai-ice', {
          peerId: 'open-ai',
          icecandidate: socketEvent.data,
        });
      }
    });

    wss.clients.forEach((client) => {
      if (client !== ws && client.readyState === WebSocket.OPEN) {
        client.send(message);
      }
    });

    ws.send(`Server received: ${message}`);
  });

  ws.on('close', () => {
    console.log('Client disconnected');
  });

  ws.on('error', (error) => {
    console.error('WebSocket error:', error);
  });
});

const io = require('socket.io')(server, {
  cors: {
    origin: '*',
    methods: ['GET', 'POST'],
  },
});

app.use(express.json());

app.use(
  cors({
    origin: '*',
    methods: ['GET', 'POST', 'PUT', 'DELETE'],
  })
);

const socketUserMapping = {};

io.on('connection', (socket) => {
  console.log('============Socket connected=============', socket.id);
  ioSockets.add(socket);
  chatServer(socket, io);

  socket.on(
    'open-ai-offer',
    async ({ peerId, sessionDescription: offer, token }) => {
      console.log('open-ai-offer', { offer, token });
      // const answer = await sendLocalDescriptionToOpenAi({ offer, token });
      // console.log('open-ai-answer', { peerId, sessionDescription: answer });
      wss.clients.forEach((client) => {
        if (client.readyState === WebSocket.OPEN) {
          const socketEvent = { event: 'offer', data: offer };
          client.send(JSON.stringify(socketEvent));
        }
      });
      // socket.emit('open-ai-answer', { peerId, sessionDescription: answer });
    }
  );

  socket.on('join', async ({ roomId, user }) => {
    console.log('============Socket join=============', {
      roomId,
      user,
    });

    // const token = await fetchEphimerialToken();
    // console.log('Token', token);

    // socket.emit('openai-session-key', { token });

    socketUserMapping[socket.id] = user;

    const clients = Array.from(io.sockets.adapter.rooms.get(roomId) || []);

    const users = clients.map((client) => {
      return socketUserMapping[client];
    });

    if (!users.includes(null)) {
      await Room.findOneAndUpdate(
        { _id: ObjectId(roomId) },
        { $set: { users: users } }
      );
    }

    clients.forEach((clientId) => {
      io.to(clientId).emit('add-peer', {
        peerId: socket.id,
        createOffer: false,
        user: user,
      });
      socket.emit('add-peer', {
        peerId: clientId,
        createOffer: true,
        user: socketUserMapping[clientId],
      });
    });

    socket.join(roomId);

    console.log('Clients connected', clients);
  });

  socket.on('relay-ice', ({ peerId, icecandidate }) => {
    console.log('============Socket relay-ice From React=============', {
      peerId,
      icecandidate,
    });

    io.to(peerId).emit('ice-candidate', {
      peerId: socket.id,
      icecandidate,
    });

    wss.clients.forEach((client) => {
      if (client.readyState === WebSocket.OPEN) {
        const socketEvent = { event: 'ice-candidates', data: icecandidate };
        client.send(JSON.stringify(socketEvent));
      }
    });
  });

  socket.on('relay-sdp', ({ peerId, sessionDescription }) => {
    io.to(peerId).emit('session-description', {
      peerId: socket.id,
      sessionDescription,
    });
  });

  socket.on('mute', ({ userId, isMuted, roomId }) => {
    const clients = Array.from(io.sockets.adapter.rooms.get(roomId) || []);
    clients.forEach((clientId) => {
      io.to(clientId).emit('mute', {
        userId,
        isMuted,
      });
    });
  });

  socket.on('un-mute', ({ userId, isMuted, roomId }) => {
    const clients = Array.from(io.sockets.adapter.rooms.get(roomId) || []);
    clients.forEach((clientId) => {
      io.to(clientId).emit('un-mute', {
        userId,
        isMuted,
      });
    });
  });

  const leaveRoom = () => {
    console.log('============Socket leave=============');
    const { rooms } = socket;

    Array.from(rooms).forEach(async (roomId) => {
      const clients = Array.from(io.sockets.adapter.rooms.get(roomId) || []);
      clients.forEach((clientId) => {
        io.to(clientId).emit('remove-peer', {
          peerId: socket.id,
          user: socketUserMapping[socket.id]?._id,
        });
      });
      socket.leave(roomId);
    });

    delete socketUserMapping[socket.id];
  };
  socket.on('leave', leaveRoom);

  socket.on('disconnecting', leaveRoom);

  socket.on('disconnect', () => {
    console.log('============Socket disconnected=============');
  });
});

db_connection();

routes(app);

app.get('/', (req, res) => {
  res.send('Welcome to sab sunno.Sab Sunno is the one of its kind application');
});

server.listen(8000, () => {
  console.log(`Server started at port ${8000}`);
});
