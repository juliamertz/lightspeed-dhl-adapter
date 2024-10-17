package main

import (
	"fmt"
	"lightspeed-dhl/config"
	"lightspeed-dhl/database"
	"lightspeed-dhl/dhl"
	"lightspeed-dhl/logger"
	"net/http"

	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := config.LoadSecrets("config.toml")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load secrets")
	}
	logger.SetupLogger(conf)

  db, err := database.Initialize("./database.db")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}


	// TODO: set up route handlers

	client := dhl.NewClient(nil)
  client.Authenticate(conf.Dhl)

	dhl.StartPolling(&client, conf, db)

	log.Info().Int("Port", *conf.Options.Port).Msg("Starting server")
	_ = http.ListenAndServe(fmt.Sprintf(":%d", *conf.Options.Port), nil)
}
