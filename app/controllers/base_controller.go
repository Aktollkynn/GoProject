package controllers

//
//import (
//	"fmt"
//	"gorm.io/driver/mysql"
//	"log"
//	"net/http"
//	"os"
//
//	"github.com/gieart87/gotoko/app/models"
//	"github.com/gieart87/gotoko/database/seeders"
//	"github.com/gorilla/mux"
//	"github.com/urfave/cli"
//	//"gorm.io/driver/mysql"
//	"gorm.io/driver/postgres"
//	"gorm.io/gorm"
//)
//
//type Server struct {
//	DB     *gorm.DB
//	Router *mux.Router
//}
//
//type AppConfig struct {
//	AppName string
//	AppEnv  string
//	AppPort string
//}
//
//type DBConfig struct {
//	DBHost     string
//	DBUser     string
//	DBPassword string
//	DBName     string
//	DBPort     string
//	DBDriver   string
//}
//
//func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
//	fmt.Println("Welcome to " + appConfig.AppName)
//
//	server.initializeDB(dbConfig)
//	server.initializeRoutes()
//}
//
//func (server *Server) Run(addr string) {
//	fmt.Printf("Listening to port %s", addr)
//	log.Fatal(http.ListenAndServe(addr, server.Router))
//}
//
//func (server *Server) initializeDB(dbConfig DBConfig) {
//	var err error
//	if dbConfig.DBDriver == "mysql" {
//		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
//		server.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	} else {
//		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
//		server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	}
//
//	if err != nil {
//		panic("Failed on connecting to the database server")
//	}
//}
//
//func (server *Server) dbMigrate() {
//	for _, model := range models.RegisterModels() {
//		err := server.DB.Debug().AutoMigrate(model.Model)
//
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//
//	fmt.Println("Database migrated successfully.")
//}
//
//func (server *Server) InitCommands(config AppConfig, dbConfig DBConfig) {
//	server.initializeDB(dbConfig)
//
//	cmdApp := cli.NewApp()
//	cmdApp.Commands = []cli.Command{
//		{
//			Name: "db:migrate",
//			Action: func(c *cli.Context) error {
//				server.dbMigrate()
//				return nil
//			},
//		},
//		{
//			Name: "db:seed",
//			Action: func(c *cli.Context) error {
//				err := seeders.DBSeed(server.DB)
//				if err != nil {
//					log.Fatal(err)
//				}
//
//				return nil
//			},
//		},
//	}
//
//	err := cmdApp.Run(os.Args)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
