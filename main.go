package main

import (
	"fmt"
	"lightspeed-dhl/config"
	"lightspeed-dhl/database"
	"lightspeed-dhl/dhl"
	"lightspeed-dhl/logger"
	"net/http"

	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := config.LoadSecrets("config.toml")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load secrets")
	}
	logger.SetupLogger(conf)

  db, err := database.Initialize("./database.db")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}


	// TODO: set up route handlers

	client := dhl.NewClient(nil)
  client.Authenticate(conf.Dhl)

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

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		log.Info().Str("Method", r.Method).Msg("Received webhook")
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

			log.Debug().Interface("Order data", orderData).Msg("Received order data from webhook")

			draft, err := dhl.WebhookToDraft(orderData, conf)
			if err != nil {
				log.Err(err).Stack().Msg("Failed to convert webhook to draft")
				return
			}

			log.Debug().Interface("Draft", draft).Msg("Transformed order data to draft")

			if !*conf.Options.DryRun {
				err, _ = dhl.CreateDraft(draft, conf)
				if err != nil {
					log.Err(err).Stack().Msg("Failed to create draft in DHL")
					return
				}

				log.Info().Str("Order reference", orderData.Order.Number).Msg("Draft created in DHL")
			}

			err = database.CreateDraft(draft.Id, draft.OrderReference, orderData.Order.Number)
			if err != nil {
				log.Err(err).Stack().Msg("Failed to create draft in database")
				return
			}

			log.Info().Str("Order reference", orderData.Order.Number).Msg("Draft created in database")

			w.WriteHeader(http.StatusOK)
		}
	})

	log.Info().Int("Port", *conf.Options.Port).Msg("Starting server")
	_ = http.ListenAndServe(fmt.Sprintf(":%d", *conf.Options.Port), nil)
}
