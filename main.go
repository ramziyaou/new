package main

import (
	"fmt"

	"os"
	"log"
	"net/http"
	"rest-go-demo/controllers"
	"rest-go-demo/database"
	"rest-go-demo/entity"
	"rest-go-demo/middleware"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql" //Required for MySQL dialect
)

func main() {
	initDB()
	// Set stop and start status of mining to false in case of incorrect exit from program earlier (e.g. program exited without stopping mining previously)
	var ss entity.StartStopCheck
	if err := database.Connector.Model(&ss).Updates(map[string]interface{}{"start": false, "stop": false}).Error; err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Starting the HTTP server on port 8080")

	router := mux.NewRouter().StrictSlash(true)
	initaliseHandlers(router)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Println(err)
		return
	}
}

func initaliseHandlers(router *mux.Router) {
	router.HandleFunc("/app/user/{id:[0-9]+}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/app/user/{id:[0-9]+}", controllers.SaveUser).Methods("POST")
	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}", controllers.GetWallet).Methods("GET")
	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}", controllers.SaveWallet).Methods("POST")
	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}/start", controllers.StartMining).Methods("OPTIONS")
	router.HandleFunc("/app/wallet/{name:[a-zA-Z]+}/stop", controllers.StopMining).Methods("OPTIONS")
	router.Use(middleware.TimerMiddleware, middleware.HTTPMethodsCheckMiddleware, middleware.AuthMiddleware)
}

func initDB() {
	user := os.Getenv("MYSQL_USER")
    pass := os.Getenv("MYSQL_PASSWORD")
    host := os.Getenv("MYSQL_HOST") 
    dbname := os.Getenv("MYSQL_DATABASE")
	config :=
		database.Config{
			ServerName: host, 
			User:       user,  
			Password:   pass,
			DB:         dbname,
		}

	connectionString := database.GetConnectionString(config)
	err := database.Connect(connectionString)
	if err != nil {
		log.Println("initDB:", err)
		return
	}
	database.Migrate(&entity.User{}, &entity.CryptoWallet{}, &entity.StartStopCheck{})
}
