package main

import "github.com/Aktollkynn/GoProject.git/app"

func main() {
	//dsn := fmt.Sprintf("host=localhost user=postgres password=online dbname=shop port=5432 sslmode=disable  TimeZone=Asia/Almaty")
	//server, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//fmt.Println(server)
	//
	//if err != nil {
	//	fmt.Printf("Failed on connecting to the database server %v", err)
	//}
	app.Run()

}
