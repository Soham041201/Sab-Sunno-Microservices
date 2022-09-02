const mongoose = require("mongoose");


const messageSchema = new mongoose.Schema({
    id: {
      type: mongoose.Schema.Types.ObjectId,
    },
    senderId: {
      type: mongoose.Schema.Types.ObjectId,
      required: true,
      ref: "Users",
    },
    recieverId: {
      type: mongoose.Schema.Types.ObjectId,
      required: true,
      ref: "Users",
    },
    messageContent:{
        type: String,
        required: true,
    },
    isDelivered: {
     type: Boolean,
     default: false
    },
    isRead: {
        type: Boolean,
        default: false
    },
    createdAt: {
      type: Date,
      default: Date.now,
    },
    updatedAt: {
      type: Date,
      default: Date.now,
    },
  });
  
  module.exports = messageSchema;
  