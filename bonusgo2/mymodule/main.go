package main

import (
	"fmt"

	"mymodule/Math"
	"mymodule/mypackage"
)

func main() {
	fmt.Println("Hello, Modules!")

	mypackage.PrintHello()
	Math.PrintMath()
}
