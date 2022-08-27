const userRoutes = require('./userRoutes')
const other_routes = require('./routes')


const routes = (app)=>{
    userRoutes(app)
    other_routes(app)
}

module.exports = routes