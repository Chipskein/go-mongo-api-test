package routes

import (
	"go-mongo-api-test/controllers/usercontroller"

	"github.com/gorilla/mux"
)

func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users", usercontroller.CreateUserEndPoint).Methods("POST")
	router.HandleFunc("/users", usercontroller.ShowAllUsersEndPoint).Methods("GET")
	router.HandleFunc("/users/{id}", usercontroller.GetUserByIdEndpoint).Methods("GET")
	return router
}
