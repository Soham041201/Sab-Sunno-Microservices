const { db_connection, Room } = require("./db-connect");
const router = require("./routes");
const express = require("express");
const app = express();
const cors = require("cors");
const { ObjectId } = require("mongodb");
require("dotenv").config();
const server = require("http").createServer(app);

const io = require("socket.io")(server, {
  cors: {
    origin: process.env.FRONT_END_URL,
    methods: ["GET", "POST"],
  },
});

app.use(express.json());

app.use(
  cors({
    origin: process.env.FRONT_END_URL,
    methods: ["GET", "POST", "PUT", "DELETE"],
  })
);

const socketUserMapping = {};

io.on("connection", (socket) => {
  console.log("============Socket connected=============", socket.id);

  socket.on("join", async ({ roomId, user }) => {
    console.log("============Socket join=============", {
      roomId,
      user,
    });

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
      io.to(clientId).emit("add-peer", {
        peerId: socket.id,
        createOffer: false,
        user: user,
      });
      socket.emit("add-peer", {
        peerId: clientId,
        createOffer: true,
        user: socketUserMapping[clientId],
      });
    });

    socket.join(roomId);
    console.log("Clients connected", clients);
  });

  socket.on("relay-ice", ({ peerId, icecandidate }) => {
    io.to(peerId).emit("ice-candidate", {
      peerId: socket.id,
      icecandidate,
    });
  });

  socket.on("relay-sdp", ({ peerId, sessionDescription }) => {
    io.to(peerId).emit("session-description", {
      peerId: socket.id,
      sessionDescription,
    });
  });

  socket.on("mute", ({ userId, isMuted, roomId }) => {
    const clients = Array.from(io.sockets.adapter.rooms.get(roomId) || []);
    clients.forEach((clientId) => {
      io.to(clientId).emit("mute", {
        userId,
        isMuted,
      });
    });
  });

  socket.on("un-mute", ({ userId, isMuted, roomId }) => {
    const clients = Array.from(io.sockets.adapter.rooms.get(roomId) || []);
    clients.forEach((clientId) => {
      io.to(clientId).emit("un-mute", {
        userId,
        isMuted,
      });
    });
  });

  const leaveRoom = () => {
    console.log("============Socket leave=============");
    const { rooms } = socket;

    Array.from(rooms).forEach(async (roomId) => {
      const clients = Array.from(io.sockets.adapter.rooms.get(roomId) || []);

      clients.forEach((clientId) => {
        io.to(clientId).emit("remove-peer", {
          peerId: socket.id,
          user: socketUserMapping[socket.id]?._id,
        });
      });
      socket.leave(roomId);
    });

    delete socketUserMapping[socket.id];
  };
  socket.on("leave", leaveRoom);

  socket.on("disconnecting", leaveRoom);
});

db_connection();

app.use(router);

app.get("/", (req, res) => {
  res.send("Welcome to sab sunno");
});

server.listen(process.env.PORT, () => {
  console.log(`Server started at port ${process.env.PORT}`);
});
