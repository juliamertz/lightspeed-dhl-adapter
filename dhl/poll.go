package dhl

import (
	"fmt"
	"lightspeed-dhl/config"
	"lightspeed-dhl/database"
	"lightspeed-dhl/lightspeed"
	"time"

	"github.com/rs/zerolog/log"
)

type Status = int8

const (
	StatusOk        Status = 0
	StatusNotFound  Status = 1
	StatusCancelled Status = 2
	StatusError     Status = 2
)

func CheckOrderStatus(client *Client, conf *config.Secrets, order *database.Order) (Status, *Label, error) {
	logger := log.With().Stack().Int("Order number", *order.LightspeedOrderId).Str("DHL draft id", *order.DhlDraftId).Logger()

	logger.Debug().Interface("order", order).Msg("Processing order")

	// Check eith DHL api if a label has been created for this order
	label, err := client.GetLabelByReference(*order.LightspeedOrderId)
	if err != nil {
		logger.Err(err).Msg("Error getting label by reference")
		return StatusError, nil, err
	}

	if label == nil {
		logger.Debug().Msg("No label found")
		return StatusNotFound, nil, nil
	}

	data, err := lightspeed.GetOrder(*order.LightspeedOrderId, conf)
	if err != nil {
		logger.Err(err).Msg("Error getting order from lightspeed, it might have been deleted")
		return StatusError, nil, err
	}

	if data.Order.Status == "cancelled" {
		// TODO: we can still return label if needed
		return StatusCancelled, nil, nil
	} else if data.Order.Status == "completed_shipped" {
		return StatusOk, label, nil
	}
	// return "processing_awaiting_shipment"

	return StatusError, nil, fmt.Errorf("Invalid order status: %v", data.Order.Status)
}

func UpdateOrderStatus(client *Client, db *database.DB, conf *config.Secrets, order *database.Order, label *Label, status Status) error {
	logger := log.With().Stack().Str("Order reference", label.OrderReference).Int("Order number", *order.LightspeedOrderId).Str("DHL draft id", *order.DhlDraftId).Logger()
	logger.Debug().Interface("Label", label).Msg("Label found")

	// Set shipment id for this order in the database
	err := db.SetShipmentId(*order.DhlDraftId, label.ShipmentId)
	if err != nil {
		logger.Err(err).Msg("Error setting shipment id")
		return err
	}

	if status == StatusCancelled {
		logger.Info().Msg("Order is cancelled, removing from database")
		// Delete cancelled orders to prevent unnecessary iterations
		err := db.DeleteDraft(*order.DhlDraftId)
		if err != nil {
			logger.Err(err).Msg("Error deleting cancelled order from database")
		}
		return err
	}

	if !*conf.Options.DryRun {
		err := lightspeed.UpdateOrderStatus(*order.LightspeedOrderId, lightspeed.UpdateOrderData{
			Status:         "completed_shipped",
			ShipmentStatus: "shipped",
		}, conf)
		if err != nil {
			logger.Err(err).Msg("Error updating order status")
			return err
		}

		logger.Debug().Msg("Order status updated to shipped")
	}

	// Set isProcessed in database
	err = db.SetProcessed(*order.DhlDraftId)
	if err != nil {
		logger.Err(err).Msg("Error setting processed")
		return err
	}
	logger.Debug().Msg("Order processed")
	return nil
}

func pollOrders(client *Client, conf *config.Secrets, db *database.DB) {
	orders, err := db.GetUnprocessed()
	if err != nil {
		log.Err(err).Stack().Msg("Failed to get all orders")
		return
	}

	log.Info().Int("Entries in database", len(orders)).Msg("Polling for labels")
	for i := range orders {
		order := orders[i]

		status, label, err := CheckOrderStatus(client, conf, &order)
		if err != nil {
			// TODO:
			panic(err)
		}

		if status != StatusOk {
			panic(status)
		}

		UpdateOrderStatus(client, db, conf, &order, label, status)
	}

}

func StartPolling(client *Client, conf *config.Secrets, db *database.DB) {
	sleepDuration := time.Duration(*conf.Options.PollingInterval) * time.Minute

	go func() {
		for {
			time.Sleep(sleepDuration)
		}
	}()
}
