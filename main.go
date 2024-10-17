package main

import (
	"fmt"
	"lightspeed-dhl/config"
	"lightspeed-dhl/database"
	"lightspeed-dhl/dhl"
	"lightspeed-dhl/logger"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := config.LoadSecrets("config.toml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load secrets")
	}

	logger.SetupLogger(conf)

	// TODO: set up route handlers

	client := dhl.NewClient(nil)
  client.Authenticate(conf.Dhl)

	dhl.StartPolling(&client, conf)

	log.Info().Int("Port", *conf.Options.Port).Msg("Starting server")
	_ = http.ListenAndServe(fmt.Sprintf(":%d", *conf.Options.Port), nil)
}
