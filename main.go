package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var mongoClient *mongo.Client

type EnviromentVariables struct {
	mongodb_url string
	port        int
}
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
}

func loadEnviroment() EnviromentVariables {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}
	mongodb_url := os.Getenv("mongodb_url")
	portString := os.Getenv("port")
	port, err := strconv.Atoi(portString)
	if err != nil {
		log.Panic("Port is not a integer")
	}

	var env EnviromentVariables
	env.mongodb_url = mongodb_url
	env.port = port

	return env
}
func createMongoClient(mongodb_url string) {
	log.Print("mongodb_url => ", mongodb_url)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb_url))
	if err != nil {
		log.Panic(err)
	}
	mongoClient = client
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createUserEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	hashedPassword, _ := HashPassword(user.Password)
	user.Password = hashedPassword
	var collection = mongoClient.Database("testdatabase").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(result)
}
func showAllUsersEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var users []User
	var collection = mongoClient.Database("testdatabase").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(users)
}
func getUserByIdEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User

	var collection = mongoClient.Database("testdatabase").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func main() {
	env := loadEnviroment()
	createMongoClient(env.mongodb_url)

	router := mux.NewRouter()
	/* ENDPOINTS */
	router.HandleFunc("/users", createUserEndPoint).Methods("POST")
	router.HandleFunc("/users", showAllUsersEndPoint).Methods("GET")
	router.HandleFunc("/users/{id}", getUserByIdEndpoint).Methods("GET")
	log.Println("Server at port:", env.port)
	http.ListenAndServe(fmt.Sprintf(":%d", env.port), router)
}
