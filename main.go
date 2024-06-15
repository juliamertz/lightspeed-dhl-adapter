package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/pkgerrors"
	"io"
	"jorismertz/lightspeed-dhl/database"
	"jorismertz/lightspeed-dhl/dhl"
	"jorismertz/lightspeed-dhl/lightspeed"
	"net/http"
	"os"

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

	database.Initialize()

	dhl.StartPolling()

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

			// err = dhl.CreateDraft(&draft)
			// if err != nil {
			// log.Err(err).Stack().Msg("Failed to create draft in DHL")
			// return
			// }

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
