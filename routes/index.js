const userRoutes = require('./userRoutes');
const other_routes = require('./roomRoutes');
const connection_routes = require('./connectionRoutes');
const message_routes = require('./messageRoutes');

const routes = (app) => {
  userRoutes(app);
  other_routes(app);
  connection_routes(app);
  message_routes(app);
};

module.exports = routes;
