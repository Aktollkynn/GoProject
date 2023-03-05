package main

import (
	"github.com/Aktollkynn/GoProject.git/app"
	"github.com/Aktollkynn/GoProject.git/app/controllers"
)

func main() {
	println("Welcome to our Project")
	println("Login http://localhost:9000/sign_in/\nHome page: http://localhost:9000/home_page/#")

	controllers.HandlerRequest()
	app.Run()

}
