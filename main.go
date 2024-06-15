package main

import (
	"encoding/json"
	"fmt"
	"io"
	"jorismertz/lightspeed-dhl/config"
	"jorismertz/lightspeed-dhl/database"
	"jorismertz/lightspeed-dhl/dhl"
	"jorismertz/lightspeed-dhl/lightspeed"
	"net/http"
	"os"

	"github.com/rs/zerolog/pkgerrors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	port = 8080
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	conf, err := config.LoadSecrets("config.toml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load secrets")
	}

	database.Initialize()

	dhl.StartPolling(conf)

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Str("Method", r.Method).Msg("Received webhook")
		if r.Method == "POST" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Err(err).Stack().Msg("Failed to read request body")
				return
			}

			var orderData lightspeed.IncomingOrder
			err = json.Unmarshal(body, &orderData)
			if err != nil {
				log.Err(err).Stack().Msg("Failed to unmarshal request body")
				return
			}

			draft, err := dhl.WebhookToDraft(orderData)
			if err != nil {
				log.Err(err).Stack().Msg("Failed to convert webhook to draft")
				return
			}

			if !*conf.Options.DryRun {
				err = dhl.CreateDraft(draft)
				if err != nil {
					log.Err(err).Stack().Msg("Failed to create draft in DHL")
					return
				}
			}

			err = database.CreateDraft(draft.Id, *draft.OrderReference, orderData.Order.Number)

			if err != nil {
				log.Err(err).Stack().Msg("Failed to create draft in database")
				return
			}

			w.WriteHeader(http.StatusOK)
		}
	})

	log.Info().Int("Port", port).Msg("Starting server")
	_ = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
