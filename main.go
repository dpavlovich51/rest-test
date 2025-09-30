package main

import (
	"my_rest_server/config"

	// Add logger
	"github.com/rs/zerolog/log"

	// Add server library
	"net/http"
)

var (
	port string = "8080"
)

func main() {

	config.SetUpLogger()
	startServer()

}

func startServer() {
	log.Info().Msgf("Starting REST server on %s port...", port)
	defer func() { log.Info().Msg("Server stoped.") }()

	router := config.WrapWithLogging(config.SetUpRouter())

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Error().Msgf("Ошибка запуска сервера: %v", err)
	}
}
