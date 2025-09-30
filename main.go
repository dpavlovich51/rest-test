package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting REST server...")
	defer func() { log.Info().Msg("Server stoped.") }()

	
}
