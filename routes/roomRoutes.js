const { ObjectId } = require("mongodb");
const { Users, Room, UserConnection } = require("../db-connect");

const router = require("express").Router();


router.post("/room", async (req, res) => {
  const { roomName, roomDescription, createdBy, users } = req.body;
  const room = await Room.create({
    roomName: roomName,
    roomDescription: roomDescription,
    createdBy: createdBy,
    users: users,
  });

  if (room) {
    return res.send({
      message: "Room created successfully",
      room: room,
      status: true,
    });
  }
  return res.send.status(400).send({
    message: "Something went wrong",
  });
});


router.get("/rooms", async (req, res) => {
  const rooms = await Room.find();
  res.send({ message: "List of rooms", rooms: rooms });
});

router.get("/room/:roomId", async (req, res) => {
  const { roomId } = req.params;
  const o_roomId = new ObjectId(roomId);
  const room = await Room.findOne({ _id: o_roomId });
  console
    .log
    // "============================Room Details==========================="
    ();
  // console.log(room);
  if (room) {
    return res.send({
      room: room,
    });
  }
  res.status(400).send({ message: "Room not found" });
});

router.get("/room/:roomId/:userId", async (req, res) => {
  const { roomId, userId } = req.params;
  const newUser = await Users.findOne({ _id: ObjectId(userId) });
  const room = await Room.findOne({ _id: ObjectId(roomId) });
  if (room.users.includes(newUser._id)) {
    return res.send({
      message: "User already in room",
      room: room,
    });
  }
  const UpdateRoom = await Room.findOneAndUpdate(
    { _id: ObjectId(roomId) },
    { $push: { users: newUser } }
  );
  if (room) {
    return res.send({
      message: "Room updated",
      room: UpdateRoom,
    });
  }
  return res.status(400).send({ message: "Something went wrong" });
});

router.delete("/room/:roomId", async (req, res) => {
  const { roomId } = req.params;
  const o_roomId = new ObjectId(roomId);
  const room = await Room.findOneAndDelete({ _id: o_roomId });

  if (room) {
    return res.send({
      message: "Room deleted",
      room: room,
    });
  }
  return res.status(400).send({ message: "Something went wrong" });
});


router.post("/user/connection/:userId/:otherUserId", async (req, res) => {
  const { userId, otherUserId } = req.params;
  const o_userId = new ObjectId(userId);
  const o_otherUserId = new ObjectId(otherUserId);

  const connection = await UserConnection.findOne({
    $or: [
      { userId: o_userId, otherUserId: o_otherUserId },
      { userId: o_otherUserId, otherUserId: o_userId },
    ],
  });

  if (connection) {
    return res.status(400).send({
      message: "User already connected",
      user: connection,
    });
  }
  const newConnection = await UserConnection.create({
    userId: o_userId,
    otherUserId: o_otherUserId,
  });
  if (newConnection) {
    return res.status(200).send({
      message: "Request sent",
      user: newConnection,
    });
  }
  return res.status(400).send({ message: "Something went wrong" });
});


  module.exports = (app)=>{
    app.use(router);
  }