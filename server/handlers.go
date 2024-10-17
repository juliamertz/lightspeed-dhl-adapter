package main

import (
	"encoding/json"
	"fmt"
	"io"
	"lightspeed-dhl/config"
	"lightspeed-dhl/database"
	"lightspeed-dhl/dhl"
	"lightspeed-dhl/lightspeed"
	"net/http"

	"github.com/rs/zerolog/log"
)

func RegisterMancoHandler(conf *config.Secrets) {
	http.HandleFunc("/stock-under-threshold", func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Str("Method", r.Method).Msg("Received request for stock under threshold")
		if r.Method == "GET" {
			data, err := lightspeed.GetStockUnderThreshold(conf)
			if err != nil {
				log.Err(err).Stack().Msg("Failed to get stock under threshold")
				return
			}

			encoded, err := json.Marshal(data)
			if err != nil {
				log.Err(err).Stack().Msg("Failed to encode stock data")
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", conf.Lightspeed.Frontend)
			w.Header().Set("Content-Type", "application/json")

			fmt.Fprintln(w, string(encoded))
		}
	})
}

func RegisterLightspeedWebhookHandler(conf *config.Secrets, client *dhl.Client) {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		log.Info().Str("Method", r.Method).Msg("Received webhook")
		if r.Method == "POST" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Err(err).Msg("Failed to read request body")
				return
			}

			var orderData lightspeed.IncomingOrder
			err = json.Unmarshal(body, &orderData)
			if err != nil {
				log.Err(err).Msg("Failed to unmarshal request body")
				return
			}

			log.Debug().Interface("Order data", orderData).Msg("Received order data from webhook")

			draft := dhl.WebhookToDraft(orderData, conf)
			log.Debug().Interface("Draft", draft).Msg("Transformed order data to draft")

			if !*conf.Options.DryRun {
				err, _ = client.CreateDraft(&draft)
				if err != nil {
					log.Err(err).Msg("Failed to create draft in DHL")
					return
				}

				log.Info().Str("Order reference", orderData.Order.Number).Msg("Draft created in DHL")
			}

			err = database.CreateDraft(draft.Id, draft.OrderReference, orderData.Order.Number)
			if err != nil {
				log.Err(err).Msg("Failed to create draft in database")
				return
			}

			log.Info().Str("Order reference", orderData.Order.Number).Msg("Draft created in database")

			w.WriteHeader(http.StatusOK)
		}
	})
}