package app

import (
	"fmt"
	"github.com/Aktollkynn/GoProject.git/database/seeders"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
}

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("Welcome to " + appConfig.AppName)

	//dsn := fmt.Sprintf("host=localhost user=postgres password=password dbname=dbname port=5432 sslmode=disable  TimeZone=Asia/Almaty")
	//server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	//var err error
	//dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Almaty", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
	//server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	////fmt.Println(server.DB)
	//
	//if err != nil {
	//  panic("Failed on connecting to the database server")
	//}
	server.initializeDB(dbConfig)
	server.InitializeRoutes()
	seeders.DBSeed(server.DB)
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (server *Server) initializeDB(dbConfig DBConfig) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Almaty", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed on connecting to the database server")
	}

	for _, model := range RegisterModels() {
		err = server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Database migrated successfully")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func Run() {

	//dsn := fmt.Sprintf("host=localhost user=postgres password=password dbname=shop port=5432 sslmode=disable  TimeZone=Asia/Almaty")
	//server, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//fmt.Println(server)
	var server = Server{}
	var appConfig = AppConfig{}
	var dbConfig = DBConfig{}

	//err := godotenv.Load()
	//if err != nil {
	//  fmt.Printf("Failed on connecting to the database server %v", err)
	//  var server = Server{}
	//  var appConfig = AppConfig{}
	//  var dbConfig = DBConfig{}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error on loading .env file")
	}

	appConfig.AppName = getEnv("APP_NAME", "TolkynSite")
	appConfig.AppEnv = getEnv("APP_ENV", "development")
	appConfig.AppPort = getEnv("APP_PORT", "9000")

	dbConfig.DBHost = getEnv("DB_HOST", "localhost")
	dbConfig.DBUser = getEnv("DB_USER", "postgres")
	dbConfig.DBPassword = getEnv("DB_PASSWORD", "online")
	dbConfig.DBName = getEnv("DB_NAME", "shop")
	dbConfig.DBPort = getEnv("DB_PORT", "5432")

	server.Initialize(appConfig, dbConfig)
	server.Run(":" + appConfig.AppPort)

}
