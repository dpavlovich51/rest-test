package config

import (
	// Add conviniet router
	"github.com/gorilla/mux"
	// Add server library
	"net/http"
	// Add logger
	"github.com/rs/zerolog/log"
)

var (
	router = mux.NewRouter()
)

func SetUpRouter() *mux.Router {

	router.HandleFunc("/messages", GetAllMessages).Methods("GET")
	router.HandleFunc("/messages/{id}", GetMessage).Methods("GET")
	router.HandleFunc("/messages", PostMessage).Methods("POST")
	router.HandleFunc("/messages/{id}", PutMessage).Methods("PUT")
	router.HandleFunc("/messages/{id}", DeleteMessage).Methods("DELETE")

	// todo: Add POST, PUT, DELETE methods
	return router
}

func GetAllMessages(response http.ResponseWriter, request *http.Request) {
	log.Info().Msg("Called GetAllMessages")
}

func GetMessage(response http.ResponseWriter, request *http.Request) {
	log.Info().Msg("Called GetMessage")
}

func PostMessage(response http.ResponseWriter, request *http.Request) {
	log.Info().Msg("Called PostMessage")
}

func PutMessage(response http.ResponseWriter, request *http.Request) {
	log.Info().Msg("Called PutMessage")
}

func DeleteMessage(response http.ResponseWriter, request *http.Request) {
	log.Info().Msg("Called DeleteMessage")
}
