package dhl

import (
	"fmt"
	"jorismertz/lightspeed-dhl/database"
	"jorismertz/lightspeed-dhl/lightspeed"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	every = 5 // minutes
)

func StartPolling() {
	go func() {
		for {
			orders, err := database.GetAll()
			if err != nil {
				log.Err(err).Msg("Failed to get all orders")
				time.Sleep(every * time.Second)
				continue
			}

			fmt.Println("polling...")
			for i := range orders {
				order := orders[i]
				fmt.Println(*order.LightspeedOrderNumber)
				label, err := GetLabelByReference(*order.LightspeedOrderNumber)
				if err != nil {
					log.Err(err).Msg("Error getting label by reference")
					continue
				}

				// This means the label has been created thus it's shipped
				if label != nil {
					err := database.SetShipmentId(*order.DhlDraftId, label.shipmentId)
					if err != nil {
						log.Err(err).Fields(order).Msg("Error setting shipment id")
						continue
					}
					// First we have to check check if the data isn't cancelled
					data, err := lightspeed.GetOrder(*order.LightspeedOrderId)
					if err != nil {
						fmt.Println(err)
					}
					if data.Order.Status == "cancelled" {
						fmt.Println("Order is cancelled, not updating status")
						continue
					} else {
						// Update status in lightspeed
						lightspeed.UpdateOrderStatus(*order.LightspeedOrderId, "shipped")
						// Set isProcessed in database
						err := database.SetProcessed(*order.DhlDraftId)
						if err != nil {
							fmt.Println(err)
						}
					}
				}

				fmt.Printf("%v", label)
			}

			time.Sleep(every * time.Second)
		}
	}()
}
