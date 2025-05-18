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

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func validateRequest(r *http.Request, secrets *config.Secrets) bool {
	validClusterId := r.Header.Get("x-cluster-id") == secrets.Lightspeed.ClusterId
	validShopId := r.Header.Get("x-shop-id") == secrets.Lightspeed.ShopId

	return validClusterId && validShopId
}

func handleGetStockUnderThreshold(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Stack().Str("method", r.Method).Str("endpoint", "/stock-under-threshold").Logger()
	logger.Debug().Msg("Getting stock under threshold")

	if r.Method == "GET" {
		data, err := lightspeed.GetStockUnderThreshold(&conf)
		if err != nil {
			logger.Err(err).Msg("Failed to get stock under threshold")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		encoded, err := json.Marshal(data)
		if err != nil {
			logger.Err(err).Msg("Failed to encode stock data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", conf.Lightspeed.Frontend)
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprintln(w, string(encoded))
	}
}

func handleLightspeedWebhook(w http.ResponseWriter, r *http.Request) {
	logger := log.With().Stack().Str("method", r.Method).Str("endpoint", "/webhoook").Logger()
	logger.Debug().Msg("Handling Lightspeed webhook")

	if r.Method == "POST" {
		valid := validateRequest(r, &conf)
		if !valid {
			logger.WithLevel(zerolog.ErrorLevel).Stack().Msg("Request cannot be verified")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Err(err).Msg("Failed to read request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var orderData lightspeed.IncomingOrder
		err = json.Unmarshal(body, &orderData)
		if err != nil {
			logger.Err(err).Msg("Failed to unmarshal request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		draft := dhl.WebhookToDraft(orderData, &conf)

		logger = logger.With().
			Str("order_reference", orderData.Order.Number).
			Interface("draft", draft).
			Logger()

		logger.Debug().Msg("Transformed order data to draft")

		if !conf.Options.DryRun {
			err, _ = dhl.CreateDraft(&draft, &conf)
			if err != nil {
				logger.Err(err).
					Interface("order_data", orderData).
					Stack().Msg("Failed to create draft in DHL")

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			logger.Info().Msg("Draft created in DHL")
		}

		unprocessedOrdersAmount.Inc()
		err = database.CreateDraft(draft.Id, draft.OrderReference, orderData.Order.Number)
		if err != nil {
			logger.
				Err(err).
				Interface("order_data", orderData).
				Stack().Msg("Failed to create draft in database")

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info().Msg("Draft created in database")

		w.WriteHeader(http.StatusOK)
	}
}
