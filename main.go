package main

import (
	"github.com/Aktollkynn/GoProject.git/app"
	"github.com/Aktollkynn/GoProject.git/app/controllers"
)

func main() {
	controllers.HandlerRequest()
	app.Run()

}
