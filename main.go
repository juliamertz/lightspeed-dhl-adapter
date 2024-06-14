package main

import (
	"encoding/json"
	"fmt"
	"io"
	"jorismertz/lightspeed-dhl/database"
	"jorismertz/lightspeed-dhl/dhl"
	"jorismertz/lightspeed-dhl/lightspeed"
	"log"
	"net/http"
	"os"
)

const (
	port = 8080
)

func main() {
	database.Initialize()
	order, err := lightspeed.GetOrder(274989157)
	if err != nil {
		panic(err)
	}
	pretty, err := json.MarshalIndent(order, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(pretty))

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

			draft := dhl.WebhookToDraft(orderData)

			// err = dhl.CreateDraft(&draft)
			// if err != nil {
			// 	panic(err)
			// }

			err = database.CreateDraft(draft.Id, *draft.OrderReference, orderData.Order.Number)

			if err != nil {
				panic(err)
			}

			w.WriteHeader(http.StatusOK)
		}
	})

	fmt.Printf("Server listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
