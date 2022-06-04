const mongoose = require("mongoose");
const userSchema = require("./models/User");
const roomSchema = require("./models/Room");
const { ServerApiVersion } = require("mongodb");

const db_connection = () => {
  mongoose.connect(process.env.MONGODB_URI, {
    useNewUrlParser: true,
    useUnifiedTopology: true,
    serverApi: ServerApiVersion.v1,
  });
  const db = mongoose.connection;
  db.on("error", console.error.bind(console, "connection error:"));
  db.once("open", function () {
    console.log("Connected to MongoDB");
  });
};

const Users = mongoose.model("Users", userSchema);
const Room = mongoose.model("Rooms", roomSchema);

module.exports = { db_connection, Users, Room };
