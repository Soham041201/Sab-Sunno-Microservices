const mongoose = require("mongoose");
const userSchema = require("./User");

const roomSchema = new mongoose.Schema({
  id: {
    type: String,
  },
  roomName: {
    type: String,
    required: true,
  },
  createdBy: {
    type: userSchema,
    required: true,
  },
  roomDescription: {
    type: String,
    required: true,
  },
  users: {
    type: [userSchema],
    required: true,
  },
});

module.exports = roomSchema;
