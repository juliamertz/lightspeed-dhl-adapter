package main

import (
	"fmt"
	"lightspeed-dhl/config"
	"lightspeed-dhl/database"
	"lightspeed-dhl/dhl"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	if os.Getenv("GO_LOG") == "debug" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if os.Getenv("ENVIRONMENT") == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

var unprocessedOrdersAmount = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "nettenshop_unprocessed_orders_amount",
		Help: "Total amount of unprocessed orders in database",
	},
)

var processedOrdersAmount = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "nettenshop_processed_orders_amount",
		Help: "Total amount of processed orders in database",
	},
)

var conf config.Secrets

func main() {
	setupLogging()

	cli := kong.Parse(&CLI)
	conf = config.LoadSecrets(cli.Args[0])

	fmt.Printf("%v\n",conf)

	database.Initialize()
	go dhl.StartPolling(&conf)

	setupPrometheus()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/stock-under-threshold", handleGetStockUnderThreshold)
	http.HandleFunc("/webhook", handleLightspeedWebhook)

	log.Info().Int("port", conf.Options.Port).Msg("Starting server")
	_ = http.ListenAndServe(fmt.Sprintf(":%d", conf.Options.Port), nil)
}

func setupPrometheus() {
	prometheus.MustRegister(unprocessedOrdersAmount)

	ordersAmount, err := database.GetUnprocessedCount()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot get unprocessed orders from database")
	}
	unprocessedOrdersAmount.Add(float64(*ordersAmount))

	prometheus.MustRegister(processedOrdersAmount)
	ordersAmount, err = database.GetProcessedCount()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot get processed orders from database")
	}
	processedOrdersAmount.Add(float64(*ordersAmount))
}
