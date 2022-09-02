const router = require("express").Router();
const { UserConnection } = require("../db-connect");
const { ObjectID } = require("bson");


router.post('/message',(req,res)=>{
    const {
        senderId,
        recieverId,
        messageContent,
        messageType,
        isDelivered,
        isRead
    } = req.body
})


module.exports = (app)=>{
    app.use(router);
  }