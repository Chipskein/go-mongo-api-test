package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Enviroment EnviromentVariables

type EnviromentVariables struct {
	MongoDBURL string
	Port       int
}

func LoadEnviroment() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}
	MongoDBURL := os.Getenv("MongoDBURL")
	portString := os.Getenv("PORT")
	Port, err := strconv.Atoi(portString)
	if err != nil {
		log.Panic("Port is not a integer")
	}

	var env EnviromentVariables
	env.MongoDBURL = MongoDBURL
	env.Port = Port
	Enviroment = env
}
