package app

import (
	"wereserve/config"

	"github.com/rs/zerolog/log"
)


func RunServer() {
	cfg := config.NewConfig()
	_, err := cfg.ConnectDB()
	if err != nil {
		log.Fatal().Msgf("Error Connection to database: %v", err)
	}

}