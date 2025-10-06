package config

import (
	// Add conviniet router
	userHandler "my_rest_server/handler"
	service "my_rest_server/service"

	"github.com/gorilla/mux"
	// Add server library
	// Add logger
)

var (
	router = mux.NewRouter()
)

func SetupRouter(userService *service.UserService) *mux.Router {
	userHandler := userHandler.NewInstance(userService)

	router.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/users", userHandler.SaveUser).Methods("POST")
	router.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	return router
}
