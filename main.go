package main

import (
	"my_rest_server/config"

	// Add logger
	"github.com/rs/zerolog/log"

	// Add server library
	"net/http"
)

var (
	port string = "8081"
)

func main() {

	startServer()

}

func startServer() {
	log.Info().Msgf("Starting REST server on %s port...", port)
	config := config.SetupApp()

	defer func() {
		log.Info().Msg("Stoping server.")
		config.Close()
	}()

	if err := http.ListenAndServe(":"+port, config.RouterHandler); err != nil {
		log.Error().Msgf("Ошибка запуска сервера: %v", err)
	}
}
