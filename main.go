package main

import (
	"github.com/Aktollkynn/GoProject.git/app"
	"github.com/Aktollkynn/GoProject.git/app/controllers"
)

func main() {
	println("Welcome to our Project")
	println("Main page  ->>   http://localhost:9000/home_page/")

	controllers.HandlerRequest()
	app.Run()

}
