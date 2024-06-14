package dhl

import (
	"fmt"
	"jorismertz/lightspeed-dhl/database"

	// "slices"

	// "jorismertz/lightspeed-dhl/database"
	"time"
)

const (
	every = 5 // minutes
)

func StartPolling() {
	go func() {
		for {
			orders, err := database.GetAll()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("polling...")
			for i := range orders {
				order := orders[i]
				fmt.Println(*order.LightspeedOrderNumber)
				label, err := GetLabelByReference(*order.LightspeedOrderNumber)
				if err != nil {
					fmt.Println(err)
				}

				// This means the label has been created thus it's shipped
				// Update status in lightspeed
				if label != nil {
					// set ShipmentId in database
					err := database.SetShipmentId(*order.DhlDraftId, label.shipmentId)
					if err != nil {
						fmt.Println(err)
					}
					// Update shipment status in lightspeed
					// First we have to check check if the data isn't cancelled
					data, err := lightspeed.GetOrder(order.LightspeedOrderId)
					if err != nil {
						fmt.Println(err)
					}
					if data.Order.Status == "cancelled" {
						fmt.Println("Order is cancelled, not updating status")
						continue
					} else {
						// Update status in lightspeed
						lightspeed.UpdateOrderStatus(order.LightspeedOrderId, "shipped")
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
