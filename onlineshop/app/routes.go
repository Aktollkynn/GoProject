package app

import (
	"github.com/Aktollkynn/GoProject.git/app/controllers"
)

func (server *Server) InitializeRoutes() {
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}
