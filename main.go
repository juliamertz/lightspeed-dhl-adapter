package main

import (
	"fmt"
	"lightspeed-dhl/config"
	"lightspeed-dhl/database"
	"lightspeed-dhl/dhl"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func main() {
  var configPath string;
  if len(os.Args) < 2 {
    configPath = "./config.toml"
  } else {
    configPath = os.Args[1]
  }

	conf, err := config.LoadSecrets(configPath)
	if err != nil {
		panic("Failed to load secrets")
	}

  db, err := database.Initialize("./database.db")
	if err != nil {
		panic("Failed to initialize database")
	}
  
	client := dhl.New(conf, dhl.DefaultCluster)
  client.Authenticate()

	SetupLogger(conf)

  RegisterMancoHandler(conf)
  RegisterLightspeedWebhookHandler(conf, &client, db)

  dhl.StartPolling(&client, conf, db)

	log.Info().Int("Port", *conf.Options.Port).Msg("Starting server")
	_ = http.ListenAndServe(fmt.Sprintf(":%d", *conf.Options.Port), nil)
}
