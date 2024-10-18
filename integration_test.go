package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lightspeed-dhl/config"
	"lightspeed-dhl/database"
	"lightspeed-dhl/dhl"
	"lightspeed-dhl/lightspeed"
	"lightspeed-dhl/logger"
	"lightspeed-dhl/server"
	"lightspeed-dhl/utils"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
)

func TestMain(t *testing.T) {
	// TODO: tailored test config
	conf, err := config.LoadSecrets("./config.toml")
	if err != nil {
		t.Fatalf("Failed to load secrets")
	}
	logger.SetupLogger(conf)

	*conf.Options.DryRun = false
	if *conf.Options.DryRun {
		t.Fatalf("no good")
	}

	db, err := database.Initialize("./tmp.db")
	if err != nil {
		t.Fatalf("Failed to initialize database")
	}

	testServerPort := randomPort(t)
	dhlMockPort := randomPort(t)

	client := dhl.Client{
		Cluster: fmt.Sprintf("http://localhost:%d", dhlMockPort),
	}
	client.Authenticate(conf.Dhl)

	var wg sync.WaitGroup
	quit := make(chan struct{})

	wg.Add(2)
	go startTestServer(&wg, quit, testServerPort, conf, &client, db)
	go startMockDhlApi(&wg, quit, dhlMockPort, conf, &client, db)

	// Test manco scraper
	testServerUrl := fmt.Sprintf("http://localhost:%d", testServerPort)

	res, err := http.Get(fmt.Sprintf("%s/stock-under-threshold", testServerUrl))
	if err != nil {
		t.Fatalf("error getting stock under threshold: %v", err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Unable to read body, error: %v", err)
	}

	var result []lightspeed.Product
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("Unable to marshal data into product, error: %v, data: %s", err, string(body))
	}

	// Test webhook
	orderData := MockLightspeedOrder()
	order, err := json.Marshal(orderData)
	if err != nil {
		t.Fatalf("Unable to marshal mock order, error: %v", err)
	}

	data := io.NopCloser(bytes.NewReader(order))
	res, err = http.Post(fmt.Sprintf("%s/webhook", client.Cluster), "application/json", data)
	if err != nil {
		t.Fatalf("error getting stock under threshold: %v", err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got: %d", res.StatusCode)
	}

	close(quit)
	wg.Wait()

	fmt.Println("Server shut down, test finished.")

	// dhl.StartPolling(&client, conf, db)

}

// TODO: Test unhappy paths
// TODO: refactor this function so it can be used for tests and normal
func startTestServer(wg *sync.WaitGroup, quit chan struct{}, port int, conf *config.Secrets, client *dhl.Client, db *database.DB) {
	defer wg.Done()

	go func() {
		server.RegisterMancoHandler(conf)
		server.RegisterLightspeedWebhookHandler(conf, client, db)

		_ = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	}()

	<-quit
}

func startMockDhlApi(wg *sync.WaitGroup, quit chan struct{}, port int, conf *config.Secrets, client *dhl.Client, db *database.DB) {
	defer wg.Done()

	go func() {

		http.HandleFunc("/authenticate/api-key", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				// TODO: Check if request is valid

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Content-Type", "application/json")

				mockAuth := dhl.AuthSession{
					AccessToken:            "fake-access-token",
					AccessTokenExpiration:  int(time.Now().Unix()) + 5,
					RefreshToken:           "fake-refresh-token",
					RefreshTokenExpiration: int(time.Now().Unix()) + 10,
				}
				encoded, err := json.Marshal(mockAuth)
				if err != nil {
					fmt.Printf("Unable to marshal mock order, error: %v", err)
					os.Exit(1)
				}

				fmt.Fprintln(w, string(encoded))
				w.WriteHeader(http.StatusOK)
			}
		})

		http.HandleFunc("/drafts", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					log.Err(err).Msg("Failed to read request body")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				var draft dhl.Draft
				err = json.Unmarshal(body, &draft)
				if err != nil {
					log.Err(err).Msg("Failed to unmarshal request body")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusCreated)
			}
		})

		_ = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	}()

	<-quit
}

func MockLightspeedOrder() lightspeed.IncomingOrder {
	order := lightspeed.Order{
		Id:             12345,
		Email:          "john.doe@example.com",
		Firstname:      "John",
		Lastname:       "Doe",
		Middlename:     "A.",
		CompanyName:    "Doe Inc.",
		Phone:          "+310634567890",
		ShipmentTitle:  "Standard Shipping",
		Number:         "ORD-98765",
		IsCompany:      true,
		Status:         "processing_awaiting_shipment",
		ShipmentStatus: "not_shipped",

		AddressBillingStreet:    "123 Main St",
		AddressBillingCity:      "Amsterdam",
		AddressBillingZipcode:   "5050 AJ",
		AddressBillingCountry:   lightspeed.CountryCode{Code: "US"},
		AddressBillingNumber:    "12",
		AddressBillingExtension: "B",

		AddressShippingStreet:    "456 Elm St",
		AddressShippingCity:      "Los Angeles",
		AddressShippingZipcode:   "5050 AJ",
		AddressShippingCountry:   lightspeed.CountryCode{Code: "US"},
		AddressShippingNumber:    "34",
		AddressShippingExtension: "C",
	}

	return lightspeed.IncomingOrder{
		Order: order,
	}
}

func randomPort(t *testing.T) (port int) {
	port, err := utils.GetFreePort()
	if err != nil {
		t.Fatalf("unable to get free port, error: %v", err)
	}

	return port
}
