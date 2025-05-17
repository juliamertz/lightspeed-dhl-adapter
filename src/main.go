package main

import (
	"fmt"
	"lightspeed-dhl/config"
	"lightspeed-dhl/database"
	"lightspeed-dhl/dhl"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var CLI struct {
	Config string `arg:"" name:"path" help:"Path to configuration file" type:"path"`
}

func setupLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if conf.Options.Debug {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// TODO: read from environment variable
	// if conf.Options.Environment == "development" {
	// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// }
}

var conf config.Secrets

func main() {
	setupLogging()

	cli := kong.Parse(&CLI)
	conf = config.LoadSecrets(cli.Args[0])

	database.Initialize()
	go dhl.StartPolling(&conf)

	http.HandleFunc("/stock-under-threshold", handleGetStockUnderThreshold)
	http.HandleFunc("/webhook", handleLightspeedWebhook)

	log.Info().Int("port", conf.Options.Port).Msg("Starting server")
	_ = http.ListenAndServe(fmt.Sprintf(":%d", conf.Options.Port), nil)
}
