package main

import (
	"fmt"
	"lightspeed-dhl/config"
	"lightspeed-dhl/database"
	"lightspeed-dhl/server"
	"lightspeed-dhl/dhl"
	"lightspeed-dhl/logger"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
)

var CLI struct {
	Config string `arg:"" name:"path" help:"Path to configuration file" type:"path"`
}

func main() {
	cli := kong.Parse(&CLI)

	conf, err := config.LoadSecrets(cli.Args[0])
	if err != nil {
		panic("Failed to load secrets")
	}

  db, err := database.Initialize("./database.db")
	if err != nil {
		panic("Failed to initialize database")
	}
  
	client := dhl.New(conf, dhl.DefaultCluster)
  client.Authenticate(conf.Dhl)

	logger.SetupLogger(conf)
	// TODO: set up route handlers

  server.RegisterMancoHandler(conf)
  server.RegisterLightspeedWebhookHandler(conf, &client, db)

  // fmt.Printf("%v", client.GetSession())
  // os.Exit(1)

  dhl.StartPolling(&client, conf, db)

	log.Info().Int("Port", *conf.Options.Port).Msg("Starting server")
	_ = http.ListenAndServe(fmt.Sprintf(":%d", *conf.Options.Port), nil)
}
