const mongoose = require("mongoose");

const userSchema = new mongoose.Schema({
  id: {
    type: String,
  },
  firstName: {
    type: String,
    required: true,
  },
  lastName: {
    type: String,
    required: true,
  },
  email: {
    type: String,
    required: true,
  },
  password: {
    type: String,
    required: true,
  },
  username: {
    type: String,
  },
  photoURL: {
    type: String,
  },
  isAuthenticated: {
    type: Boolean,
    required: true,
  },
  about:{
    type: String
  }
});

module.exports = userSchema;
