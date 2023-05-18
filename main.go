// `main.go`
package main

import (
	"github.com/Aktollkynn/GoProject.git/app/controllers"
	// "flag"
)

func main() {

	println("~ Welcome!")
	println("~ http://localhost:8000/welcome/")
	controllers.HandlerRequest()

}
