const { ObjectId } = require("mongodb");
const { Users, Room } = require("./db-connect");

const router = require("express").Router();

router.post("/register", async (req, res) => {
  console.log(req.body);
  const { firstName, lastName, email, password, username, photoURL } = req.body;

  const user = await Users.findOne({ email: email });
  // console.log("=======================User Data======================");
  // console.log(user);
  if (user === null) {
    const user = await Users.create({
      firstName: firstName,
      lastName: lastName,
      email: email,
      password: password,
      username: username,
      photoURL: photoURL,
    });
    if (user) {
      return res.send({
        message: "User created successfully",
        user: user,
      });
    }
  }

  return res.status(400).send({
    message: "User already exists",
    user: user,
  });
});

router.post("/room", async (req, res) => {
  const { roomName, roomDescription, createdBy, users } = req.body;
  const room = await Room.create({
    roomName: roomName,
    roomDescription: roomDescription,
    createdBy: createdBy,
    users: users,
  });

  // console.log("============================New Room Details===========================")
  // console.log(room);

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

router.post("/user", async (req, res) => {
  const { email, id } = req.body;
  const o_id = new ObjectId(id);
  const user = await Users.findOne({ $or: [{ email: email }, { _id: o_id }] });
  if (user) {
    return res.send({
      user: user,
    });
  }
  res.status(400).send("User not found");
});

router.get("/user/:userId", async (req, res) => {
  const { userId } = req.params;
  const o_userId = new ObjectId(userId);
  const user = await Users.findOne({ _id: o_userId });
  if (user) {
    return res.status(200).send({
      user: user,
    });
  }
  res.status(400).send({ message: "User not found" });
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
    // console.log("============================Updated Room Details===========================")
    // console.log(UpdateRoom);
    return res.send({
      message: "Room updated",
      room: room,
    });
  }
  return res.status(400).send({ message: "Something went wrong" });
});

router.delete("/room/:roomId", (req, res) => {
  const { roomId } = req.params;
  const o_roomId = new ObjectId(roomId);
  const room = Room.findOneAndDelete({ _id: o_roomId });
  if (room) {
    return res.send({
      message: "Room deleted",
      room: room,
    });
  }
  return res.status(400).send({ message: "Something went wrong" });
});

module.exports = router;
