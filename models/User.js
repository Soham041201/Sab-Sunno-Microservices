const mongoose = require("mongoose");

const userSchema = new mongoose.Schema({
  id: {
    type: String,
  },
  firstName: {
    type: String,
  },
  lastName: {
    type: String,
  },
  email: {
    type: String,
  },
  password: {
    type: String,
  },
  username: {
    type: String,
  },
  photoURL: {
    type: String,
  },
  isAuthenticated: {
    type: Boolean,
  },
  about:{
    type: String
  }
});

module.exports = userSchema;
