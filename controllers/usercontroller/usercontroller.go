package usercontroller

import (
	"context"
	"encoding/json"
	"go-mongo-api-test/utils/database"
	"go-mongo-api-test/utils/password"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
}

func CreateUserEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	hashedPassword, _ := password.HashPassword(user.Password)
	user.Password = hashedPassword
	var mongoClient = database.GetMongoClient()
	var collection = mongoClient.Database("testdatabase").Collection("users")
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(result)
}
func ShowAllUsersEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var users []User
	var mongoClient = database.GetMongoClient()
	var collection = mongoClient.Database("testdatabase").Collection("users")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
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
func GetUserByIdEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User
	var mongoClient = database.GetMongoClient()
	var collection = mongoClient.Database("testdatabase").Collection("users")
	err := collection.FindOne(context.TODO(), User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(user)
}
func UpdateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var user User
	var user2 User
	json.NewDecoder(request.Body).Decode(&user)

	var mongoClient = database.GetMongoClient()
	var collection = mongoClient.Database("testdatabase").Collection("users")
	err := collection.FindOne(context.TODO(), User{ID: user.ID}).Decode(&user2)

	if user.Password != "" {
		hashedPassword, _ := password.HashPassword(user.Password)
		user2.Password = hashedPassword
	}
	if user.Name != "" {
		user2.Name = user.Name
	}
	if user.Email != "" {
		user2.Email = user.Email
	}

	filter := bson.D{primitive.E{Key: "_id", Value: user2.ID}}
	update := bson.D{
		primitive.E{Key: "$set",
			Value: bson.D{
				primitive.E{Key: "email", Value: user2.Email},
				primitive.E{Key: "password", Value: user2.Password},
				primitive.E{Key: "name", Value: user2.Name},
			},
		}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(result)
}
func DeleteUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var mongoClient = database.GetMongoClient()
	var collection = mongoClient.Database("testdatabase").Collection("users")
	_, err := collection.DeleteOne(context.TODO(), User{ID: id})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write([]byte(`{"message":"05_DELETED"}`))

}
