package dhl

import (
	"lightspeed-dhl/config"
	"lightspeed-dhl/database"
	"lightspeed-dhl/lightspeed"
	"time"

	"github.com/rs/zerolog/log"
)

func StartPolling(conf *config.Secrets) {
	sleepDuration := time.Duration(*conf.Options.PollingInterval) * time.Minute
	go func() {
		for {
			orders, err := database.GetUnprocessed()
			if err != nil {
				log.Err(err).Stack().Msg("Failed to get all orders")
				time.Sleep(sleepDuration)
				continue
			}

			log.Info().Int("Entries in database", len(orders)).Msg("Polling for labels")
			for i := range orders {
				order := orders[i]

				logger := log.With().Stack().Str("Order number", *order.LightspeedOrderNumber).Logger()
				logger.Debug().Interface("order", order).Msg("Processing order")

				label, err := GetLabelByReference(*order.LightspeedOrderNumber)
				if err != nil {
					logger.Err(err).Msg("Error getting label by reference")
					continue
				}

				// If a label has been created for our order we can assume this means it has been shipped by DHL.
				// We check for a cancelled state in lightspeed to make sure of this
				if label != nil {
					logger.Debug().Interface("Label", label).Msg("Label found")
					err := database.SetShipmentId(*order.DhlDraftId, label.shipmentId)
					if err != nil {
						logger.Err(err).Msg("Error setting shipment id")
						continue
					}

					data, err := lightspeed.GetOrder(*order.LightspeedOrderId)
					if err != nil {
						logger.Err(err).Msg("Error getting order")
						continue
					}

					isCancelled := data.Order.Status == "cancelled"
					logger.Debug().Bool("status", isCancelled).Msg("Order cancelled status")

					if isCancelled {
						logger.Info().Msg("Order is cancelled, not updating status")
						continue
					}

					if !*conf.Options.DryRun {
						err := lightspeed.UpdateOrderStatus(*order.LightspeedOrderId, "shipped")
						if err != nil {
							logger.Err(err).Msg("Error updating order status")
							continue
						}

						logger.Debug().Msg("Order status updated to shipped")
					}
					// Set isProcessed in database
					err = database.SetProcessed(*order.DhlDraftId)
					if err != nil {
						logger.Err(err).Msg("Error setting processed")
						continue
					}
					logger.Debug().Msg("Order processed")
				}
			}

			time.Sleep(sleepDuration)
		}
	}()
}
