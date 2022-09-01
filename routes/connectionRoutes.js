const router = require("express").Router();
const { UserConnection } = require("../db-connect");
const { ObjectID } = require("bson");


router.post("/user/connection", async (req, res) => {
    const { userId, otherUserId } = req.body;
    const o_userId = new ObjectID(userId);
    const o_otherUserId = new ObjectID(otherUserId);
  
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
  

  router.post("/connection/status", async (req, res) => {
    const { userId, otherUserId } = req.body;
    const o_userId = new ObjectID(userId);
    const o_otherUserId = new ObjectID(otherUserId);
  
    const connection = await UserConnection.findOne({
      $or: [
        { userId: o_userId, otherUserId: o_otherUserId },
        { userId: o_otherUserId, otherUserId: o_userId },
      ],
    });
  
    if (connection) {
      return res.status(200).send({
        message: "Connection Exists",
        user: connection,
      });
    }
    return res.status(400).send({ message: "Connection not present" });
  });


module.exports = (app)=>{
    app.use(router);
  }