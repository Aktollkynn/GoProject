package main

import (
	"github.com/Aktollkynn/GoProject.git/app"
	"github.com/Aktollkynn/GoProject.git/app/controllers"
)

func main() {
	println("Welcome to our Project")
	println("http://localhost:9000/sign_in/")

	controllers.HandlerRequest()
	app.Run()

}
