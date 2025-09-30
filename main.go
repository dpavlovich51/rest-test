package main

import (
	// Add logger
	"github.com/rs/zerolog/log"
	// Add conviniet router
	"github.com/gorilla/mux"
	// Add server library
	"net/http"
)

var (
	port string = "8080"
)

func main() {
	startServer()

}

func startServer() {
	log.Info().Msgf("Starting REST server on %s port...", port)
	defer func() { log.Info().Msg("Server stoped.") }()

	router := mux.NewRouter()

	router.HandleFunc("/messages", GetAllMessages).Methods("GET")
	router.HandleFunc("/messages/{id}", GetMessage).Methods("GET")
	// todo: Add POST, PUT, DELETE methods

	if err := http.ListenAndServe(":" + port, router); err != nil {
		log.Error().Msgf("Ошибка запуска сервера: %v", err)
	}
}

func GetAllMessages(response http.ResponseWriter, request *http.Request) {
	log.Info().Msg("Called GetAllMessages")
}

func GetMessage(response http.ResponseWriter, request *http.Request) {
	log.Info().Msg("Called GetMessage")
}
