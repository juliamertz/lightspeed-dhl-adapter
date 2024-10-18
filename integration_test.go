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
)

type TestServer struct {
	t    *testing.T
	wg   *sync.WaitGroup
	quit chan struct{}
	port int
}

func (s *TestServer) close() {
	close(s.quit)
	s.wg.Wait()
}

func newServer(t *testing.T) TestServer {
	var wg sync.WaitGroup
	wg.Add(1)
	return TestServer{
		wg:   &wg,
		quit: make(chan struct{}),
		port: randomPort(t),
		t:    t,
	}
}

func TestMain(t *testing.T) {
	// gofakeit.Seed(8675309)
	// var mockDraft dhl.Draft
	// gofakeit.Struct(&mockDraft)
	// t.Fatalf("mockdraft: %v", mockDraft)

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
	t.Cleanup(func() { os.Remove("./tmp.db") })

	testServer := newServer(t)
	dhlMockServer := newServer(t)

	client := dhl.New(conf, fmt.Sprintf("http://localhost:%d", dhlMockServer.port))
	client.Authenticate(conf.Dhl)

	go startTestServer(&testServer, conf, &client, db)
	go startMockDhlApi(&dhlMockServer)

	// Test manco scraper
	testServerUrl := fmt.Sprintf("http://localhost:%d", testServer.port)

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

	_, err = db.GetDraft(orderData.Order.Id)
	if err != nil {
		t.Fatalf("draft not created in database, error: %v", err)
	}

	testServer.close()

	// // Test order status updating
	// orders, err := db.GetUnprocessed()
	// if err != nil {
	// 	t.Fatalf("Failed to get unprocessed orders, error: %v", err)
	// 	return
	// }
	
  // TODO: find meaningful way to test this
	// log.Info().Int("Entries in database", len(orders)).Msg("Polling for labels")
	// for i := range orders {
	// 	order := orders[i]
	// 	// "cancelled" / "completed_shipped" / "processing_awaiting_shipment"
	// 	status, label, err := dhl.CheckOrderStatus(&client, conf, &order)
	// 	if err != nil {
	// 		t.Fatalf("Unable to get order status, error: %v", err)
	// 	}
	//
	// 	if status != dhl.StatusOk {
	// 		t.Fatalf("Expected status %v got: %v", dhl.StatusOk, status)
	// 	}
	//
	// 	err = dhl.UpdateOrderStatus(&client, db, conf, &order, label, status)
	// 	if status != dhl.StatusOk {
	// 		t.Fatalf("Unable to update order status, error: %v", err)
	// 	}
	// }

	dhlMockServer.close()
}

// TODO: Test unhappy paths
func startTestServer(s *TestServer, conf *config.Secrets, client *dhl.Client, db *database.DB) {
	defer s.wg.Done()

	go func() {
		server.RegisterMancoHandler(conf)
		server.RegisterLightspeedWebhookHandler(conf, client, db)

		_ = http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)

	}()

	<-s.quit
}

func startMockDhlApi(s *TestServer) {
	defer s.wg.Done()

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
					s.t.Fatalf("Unable to marshal mock order, error: %v", err)
				}

				fmt.Fprintln(w, string(encoded))
				w.WriteHeader(http.StatusOK)
			}
		})

		http.HandleFunc("/drafts", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					s.t.Fatalf("Failed to read request body")
				}

				var draft dhl.Draft
				err = json.Unmarshal(body, &draft)
				if err != nil {
					s.t.Fatalf("Failed to unmarshal request body")
				}

				w.WriteHeader(http.StatusCreated)
			}
		})

		// http.HandleFunc("/labels", func(w http.ResponseWriter, r *http.Request) {
		// 	if r.Method == "GET" {
		// 		// TODO: get some data
		// 		labels := []dhl.Label{
		// 			{},
		// 		}
		// 		ser, err := json.Marshal(labels)
		// 		if err != nil {
		// 			log.Err(err).Msg("Failed to marshal response body")
		// 			w.WriteHeader(http.StatusInternalServerError)
		// 			return
		// 		}
		//
		// 		w.Header().Set("Content-Type", "application/json")
		// 		fmt.Fprintln(w, string(ser))
		// 		w.WriteHeader(http.StatusOK)
		// 	}
		// })

		_ = http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)

	}()

	<-s.quit
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
