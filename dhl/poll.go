package dhl

import (
	"fmt"
	"jorismertz/lightspeed-dhl/config"
	"jorismertz/lightspeed-dhl/database"
	"jorismertz/lightspeed-dhl/lightspeed"
	"time"

	"github.com/rs/zerolog/log"
)

func StartPolling(conf *config.Secrets) {
	sleepDuration := time.Duration(*conf.Options.PollingInterval) * time.Minute
	go func() {
		for {
			orders, err := database.GetAll()
			if err != nil {
				log.Err(err).Stack().Msg("Failed to get all orders")
				time.Sleep(sleepDuration)
				continue
			}

			fmt.Println("polling...")
			for i := range orders {
				order := orders[i]

				label, err := GetLabelByReference(*order.LightspeedOrderNumber)
				if err != nil {
					log.Err(err).Stack().Msg("Error getting label by reference")
					continue
				}

				// If a label has been created for our order we can assume this means it has been shipped by DHL.
				// We check for a cancelled state in lightspeed to make sure of this
				if label != nil {
					err := database.SetShipmentId(*order.DhlDraftId, label.shipmentId)
					if err != nil {
						log.Err(err).Stack().Fields(order).Msg("Error setting shipment id")
						continue
					}

					data, err := lightspeed.GetOrder(*order.LightspeedOrderId)
					if err != nil {
						log.Err(err).Stack().Fields(order).Msg("Error getting order")
						continue
					}

					if data.Order.Status == "cancelled" {
						log.Info().Str("Order reference", *order.LightspeedOrderNumber).Msg("Order is cancelled, not updating status")
						continue
					} else {
						// Update status in lightspeed
						if !*conf.Options.DryRun {
							err := lightspeed.UpdateOrderStatus(*order.LightspeedOrderId, "shipped")
							if err != nil {
								log.Err(err).Stack().Fields(order).Msg("Error updating order status")
								continue
							}
						}
						// Set isProcessed in database
						err := database.SetProcessed(*order.DhlDraftId)
						if err != nil {
							log.Err(err).Stack().Fields(order).Msg("Error setting processed")
							continue
						}
					}
				}
			}

			time.Sleep(sleepDuration)
		}
	}()
}
