package main

import (
	"fmt"
	"go-mongo-api-test/routes"
	"go-mongo-api-test/utils/database"
	"go-mongo-api-test/utils/env"
	"log"
	"net/http"
)

func main() {
	env.LoadEnviroment()
	database.CreateMongoClient(env.Enviroment.MongoDBURL)
	router := routes.CreateRouter()
	log.Println("Server at port:", env.Enviroment.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", env.Enviroment.Port), router)
}
