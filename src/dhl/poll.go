package dhl

import (
	"lightspeed-dhl/config"
	"lightspeed-dhl/database"
	"lightspeed-dhl/lightspeed"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

func poll(conf *config.Secrets, processedCount prometheus.Gauge, unprocessedCount prometheus.Gauge) error {
	orders, err := database.GetUnprocessed()
	if err != nil {
		log.Err(err).Stack().Msg("Failed to get all orders")
		return err
	}

	log.Info().Int("entries_in_database", len(orders)).Msg("Polling for labels")

	for i := range orders {
		order := orders[i]

		// Create base of logging context
		baseLogger := log.With().Stack().Int("order_number", *order.LightspeedOrderId).Str("draft_id", *order.DhlDraftId)
		logger := baseLogger.Logger()

		logger.Debug().Interface("order", order).Msg("Processing order")

		// Check with DHL api if a label has been created for this order
		label, err := GetLabelByReference(*order.LightspeedOrderId, conf, 0)
		if err != nil {
			logger.Err(err).Msg("Error getting label by reference")
			continue
		}

		if label == nil {
			logger.Debug().Msg("No label found")
			continue
		}

		// Add extra information to logging context
		logger = baseLogger.Str("order_reference", label.OrderReference).Logger()

		logger.Debug().Interface("label", label).Msg("Label found")
		// Set shipment id for this order in the database
		err = database.SetShipmentId(*order.DhlDraftId, label.ShipmentId)
		if err != nil {
			logger.Err(err).Msg("Error setting shipment id")
			continue
		}

		data, err := lightspeed.GetOrder(*order.LightspeedOrderId, conf)
		if err != nil {
			logger.Err(err).Msg("Error getting order from lightspeed, it might have been deleted")
			continue
		}

		isCancelled := data.Order.Status == "cancelled"
		logger.Debug().Bool("status", isCancelled).Msg("Order cancelled status")

		if isCancelled {
			logger.Info().Msg("Order is cancelled, removing from database")
			// Delete cancelled orders to prevent unnecessary iterations
			err := database.DeleteDraft(*order.DhlDraftId)
			if err != nil {
				logger.Err(err).Msg("Error deleting cancelled order from database")
			}
			continue
		}

		if !conf.Options.DryRun {
			err := lightspeed.UpdateOrderStatus(*order.LightspeedOrderId, lightspeed.UpdateOrderData{
				Status:         "completed_shipped",
				ShipmentStatus: "shipped",
			}, conf)
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
		processedCount.Inc()
		unprocessedCount.Dec()
		logger.Debug().Msg("Order processed")
	}

	log.Info().Msg("Done polling for labels")

	return nil
}

func StartPolling(
	conf *config.Secrets,
	pollingDuration prometheus.Histogram,
	processedCount prometheus.Gauge,
	unprocessedCount prometheus.Gauge,
) {
	sleepDuration := time.Duration(conf.Options.PollingInterval) * time.Minute

	for {
		start := time.Now()

		err := poll(conf, processedCount, unprocessedCount)
		if err != nil {
			log.Err(err).Msg("polling failed")
		}

		duration := time.Since(start).Seconds()
		pollingDuration.Observe(duration)

		time.Sleep(sleepDuration)
	}
}
