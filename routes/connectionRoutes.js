const router = require("express").Router();
const {UserConnection} = require('../models/UserConnection')



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