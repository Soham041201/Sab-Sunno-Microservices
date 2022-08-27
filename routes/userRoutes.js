const router = require("express").Router();
const { ObjectID } = require("bson");
const { Users} = require("../db-connect");




router.get("/users", async (req, res) => {
  const users = await Users.find({});
  res.send(JSON.stringify({ message: "List of users", users: users }));
});

router.post("/field/:userId", async (req, res) => {
    console.log("=======================User Data======================");
    const { userId } = req.params;
    const { field, value } = req.body;
    console.log(req.body);
    console.log(userId)
    const o_userId = new ObjectID(userId);
    const user = await Users.findByIdAndUpdate(
      o_userId,
      {
        $set: {
          [field]: value,
        },
      }
    );
    if (user) {
      console.log(user);
      return res.status(200).send({
        message: "User updated",
        user: user,
      });
    }
    res.status(400).send({ message: "User not found" });
  });


  module.exports = (app)=>{
    app.use(router);
  }