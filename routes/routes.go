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
	router.HandleFunc("/users/{id}/delete", usercontroller.DeleteUserEndpoint).Methods("GET")
	router.HandleFunc("/users/update", usercontroller.UpdateUserEndpoint).Methods("POST")
	//handlers.LoggingHandler(os.Stdout, router)
	return router
}
