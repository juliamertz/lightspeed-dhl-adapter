package main

import (
	"encoding/json"
	"fmt"
	"io"
	"jorismertz/lightspeed-dhl/database"
	"jorismertz/lightspeed-dhl/dhl"
	"jorismertz/lightspeed-dhl/lightspeed"
	"jorismertz/lightspeed-dhl/secrets"
	"log"
	"net/http"
	"os"
)

const (
	port = 8080
)

func main() {
	database.Initialize()

	config, err := secrets.LoadSecrets("config.toml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", config)

	os.Exit(1)
	dhl.StartPolling()

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			var orderData lightspeed.IncomingOrder
			err = json.Unmarshal(body, &orderData)
			if err != nil {
				panic(err)
			}

			draft := lightspeed.WebhookToDraft(orderData)
			// err = dhl.CreateDraft(&draft, secrets.DhlCredentials())
			// if err != nil {
			// 	panic(err)
			// }

			err = database.CreateDraft(&draft)
			if err != nil {
				panic(err)
			}

			// prettyDraft, _ := json.MarshalIndent(draft, "", "  ")
			// fmt.Print(string(prettyDraft))
		}
	})

	fmt.Printf("Server listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
