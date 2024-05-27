const userRoutes = require('./userRoutes');
const other_routes = require('./roomRoutes');
const connection_routes = require('./connectionRoutes');

const routes = (app) => {
  userRoutes(app);
  other_routes(app);
  connection_routes(app);
};

module.exports = routes;
