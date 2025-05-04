package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"lightspeed-dhl/config"
	"lightspeed-dhl/database"
	"lightspeed-dhl/dhl"
	"lightspeed-dhl/lightspeed"
	"lightspeed-dhl/logger"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
)

var CLI struct {
	Config string `arg:"" name:"path" help:"Path to configuration file" type:"path"`
}

func validateSignature(r *http.Request, apiSecret string) (bool, error) {
	signature := r.Header.Get("x-signature")
	if signature == "" {
		return false, fmt.Errorf("missing x-signature header")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return false, err
	}

	data := append(body, []byte(apiSecret)...)
	hash := md5.Sum(data)
	calculatedSig := hex.EncodeToString(hash[:])

	return calculatedSig == signature, nil
}

func main() {
	cli := kong.Parse(&CLI)
	conf, err := config.LoadSecrets(cli.Args[0])
	if err != nil {
		fmt.Println("Error loading secrets: ", err)
		os.Exit(1)
	}

	logger.SetupLogger(conf)
	database.Initialize()

	dhl.StartPolling(conf)

	http.HandleFunc("/stock-under-threshold", func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Str("Method", r.Method).Msg("Received request for stock under threshold")
		if r.Method == "GET" {
			data, err := lightspeed.GetStockUnderThreshold(conf)
			if err != nil {
				log.Err(err).Stack().Msg("Failed to get stock under threshold")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			encoded, err := json.Marshal(data)
			if err != nil {
				log.Err(err).Stack().Msg("Failed to encode stock data")
				w.WriteHeader(http.StatusInternalServerError)
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
			valid, err := validateSignature(r, conf.Lightspeed.Key)
			if err != nil {
				log.Err(err).Stack().Msg("Failed to verify request signature")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if !valid {
				log.Err(err).Stack().Msg("Request signature is invalid")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Err(err).Stack().Msg("Failed to read request body")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			var orderData lightspeed.IncomingOrder
			err = json.Unmarshal(body, &orderData)
			if err != nil {
				log.Err(err).Stack().Msg("Failed to unmarshal request body")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			log.Debug().Interface("Order data", orderData).Msg("Received order data from webhook")

			draft := dhl.WebhookToDraft(orderData, conf)
			log.Debug().Interface("Draft", draft).Msg("Transformed order data to draft")

			if !*conf.Options.DryRun {
				err, _ = dhl.CreateDraft(&draft, conf)
				if err != nil {
					log.Err(err).Stack().Msg("Failed to create draft in DHL")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				log.Info().Str("Order reference", orderData.Order.Number).Msg("Draft created in DHL")
			}

			err = database.CreateDraft(draft.Id, draft.OrderReference, orderData.Order.Number)
			if err != nil {
				log.Err(err).Stack().Msg("Failed to create draft in database")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Info().Str("Order reference", orderData.Order.Number).Msg("Draft created in database")

			w.WriteHeader(http.StatusOK)
		}
	})

	log.Info().Int("Port", *conf.Options.Port).Msg("Starting server")
	_ = http.ListenAndServe(fmt.Sprintf(":%d", *conf.Options.Port), nil)
}
