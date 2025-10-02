package config

import (
	// Add conviniet router
	"fmt"
	userHandler "my_rest_server/handler"
	"my_rest_server/storage"

	"github.com/gorilla/mux"
	// Add server library
	// Add logger
)

var (
	router = mux.NewRouter()
)

func SetUpRouter() *mux.Router {
	userHandler := newUserHandler()

	router.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/users", userHandler.SaveUser).Methods("POST")
	router.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	return router
}

func newUserHandler() *userHandler.UserHandler {
	storage, err := storage.NewClient("localhost:6379", "", 0)
	if err != nil {
		panic(fmt.Errorf("failed to connect to storage. error: %s", err))
	}
	userHandler := userHandler.NewInstance(storage)
	return userHandler
}
