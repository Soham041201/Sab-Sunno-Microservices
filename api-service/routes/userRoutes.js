const router = require('express').Router();
const { ObjectID } = require('bson');
const { Users } = require('../db-connect');

router.get('/users/:userId', async (req, res) => {
  const { userId } = req.params;
  console.log(userId);
  let users = await Users.find({});
  users = users.filter(
    (user) => user._id.toString() != userId && user.photoURL != ''
  );
  res.send(JSON.stringify({ message: 'List of users', users: users }));
});

router.post('/field/:userId', async (req, res) => {
  console.log('=======================User Data======================');
  const { userId } = req.params;
  const { field, value } = req.body;
  console.log(req.body);
  console.log(userId);
  const o_userId = new ObjectID(userId);
  const user = await Users.findByIdAndUpdate(o_userId, {
    $set: {
      [field]: value,
    },
  });
  if (user) {
    console.log(user);
    return res.status(200).send({
      message: 'User updated',
      user: user,
    });
  }
  res.status(400).send({ message: 'User not found' });
});

router.post('/register', async (req, res) => {
  console.log(req.body.phoneNumber);
  const {
    firstName,
    lastName,
    email,
    password,
    username,
    photoURL,
    phoneNumber,
    isAuthenticated,
  } = req.body;

  const user = await Users.findOne({
    $or: [{ email: email }],
  });
  // console.log("=======================User Data======================");
  console.log(user);
  if (user === null) {
    const user = await Users.create({
      firstName: firstName,
      lastName: lastName,
      email: email,
      password: password,
      username: username,
      photoURL: photoURL,
      isAuthenticated: false,
      phoneNumber: phoneNumber,
    });
    if (user) {
      return res.send({
        message: 'User created successfully',
        user: user,
      });
    }
  }

  return res.status(400).send({
    message: 'User already exists',
    user: user,
  });
});

router.post('/login', async (req, res) => {
  const { email, password } = req.body;
  console.log(res.body);
  const user = await Users.findOne({
    email: email,
  });

  if (user?.password == password) {
    return res.status(200).send({
      message: 'Logged in successfully',
      user: user,
    });
  } else if (user.password != password) {
    return res.status(400).send({
      message: 'Password in incorrect',
      user: user,
    });
  }
  return res.send(400).send({
    message: 'Something went wrong',
  });
});

router.post('/user', async (req, res) => {
  const { email, id, phoneNumber } = req.body;
  console.log(req.body);
  const o_id = new ObjectId(id);
  const user = await Users.findOne({
    $or: [{ email: email }, { _id: o_id }, { phoneNumber: phoneNumber }],
  });
  if (user) {
    console.log(user);
    return res.send({
      user: user,
    });
  }
  res.status(400).send('User not found');
});

router.get('/user/:userId', async (req, res) => {
  const { userId } = req.params;
  const o_userId = new ObjectID(userId.trim());
  const user = await Users.findById(o_userId);
  if (user) {
    console.log(user);
    return res.status(200).json({
      message: 'Data found',
      user: user,
    });
  }
  res.status(400).send({ message: 'User not found' });
});

router.put('/user/:userId', async (req, res) => {
  console.log('=======================User field Data======================');

  const { userId } = req.params;
  const { username, photoURL } = req.body;
  const o_userId = new ObjectID(userId);
  const user = await Users.findOneAndUpdate(
    { _id: o_userId },
    {
      $set: {
        username: '@' + username,
        photoURL: photoURL,
        isAuthenticated: true,
      },
    }
  );
  if (user) {
    return res.status(200).send({
      message: 'User updated',
      user: user,
    });
  }
  res.status(400).send({ message: 'User not found' });
});

router.put('/user/update/:userId', async (req, res) => {
  const { userId } = req.params;
  const user = req.body;
  console.log('user', user);
  const o_userId = new ObjectID(userId);
  console.log(o_userId);
  const userD = await Users.findOneAndUpdate({ _id: o_userId }, user);
  if (userD) {
    console.log(userD);
    return res.status(200).send({
      message: 'User updated',
      user: userD,
    });
  }
  res.status(400).send({ message: 'User not found' });
});

module.exports = (app) => {
  app.use(router);
};
