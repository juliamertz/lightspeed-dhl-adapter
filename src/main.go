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

var (
	unprocessedOrdersAmount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nettenshop_unprocessed_orders_amount",
			Help: "Total amount of unprocessed orders in database",
		},
	)
	processedOrdersAmount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "nettenshop_processed_orders_amount",
			Help: "Total amount of processed orders in database",
		},
	)
	pollingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "nettenshop_polling_duration",
			Help:    "Histogram of duration polling DHL for labels took.",
			Buckets: prometheus.DefBuckets,
		},
	)
)

var conf config.Secrets

func main() {
	setupLogging()

	cli := kong.Parse(&CLI)
	conf = config.LoadSecrets(cli.Args[0])

	database.Initialize()
	go dhl.StartPolling(&conf, pollingDuration, processedOrdersAmount, unprocessedOrdersAmount)

	setupPrometheus()
	serve()
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

func setupPrometheus() {
	prometheus.MustRegister(unprocessedOrdersAmount)
	prometheus.MustRegister(processedOrdersAmount)
	prometheus.MustRegister(pollingDuration)

	ordersAmount, err := database.GetUnprocessedCount()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot get unprocessed orders from database")
	}
	unprocessedOrdersAmount.Add(float64(*ordersAmount))

	ordersAmount, err = database.GetProcessedCount()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot get processed orders from database")
	}
	processedOrdersAmount.Add(float64(*ordersAmount))

}

func serve() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/stock-under-threshold", handleGetStockUnderThreshold)
	http.HandleFunc("/webhook", handleLightspeedWebhook)

	log.Info().Int("port", conf.Options.Port).Msg("Starting server")
	_ = http.ListenAndServe(fmt.Sprintf(":%d", conf.Options.Port), nil)
}
