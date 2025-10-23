use anyhow::{Context, Result};
use futures::StreamExt;
use tracing::{debug, error, info, warn};

use crate::{
    AdapterState,
    database::{self, models::Order},
    utils::ToUuid,
};
use lightspeed::{OrderStatus, ShipmentStatus};

async fn reconcile_order_status(state: &AdapterState, order: &Order) -> Result<OrderStatus> {
    let order_id = order.lightspeed_order_id as u64;
    let draft_id = order
        .dhl_draft_id
        .context("order has been created without a draft id")?;

    database::incr_poll_count(&state.pool, order.lightspeed_order_id).await?;

    if let Some(label) = state.dhl.get_label(order_id).await? {
        info!({ shipment_id = label.shipment_id }, "found label for order");

        let data = state
            .lightspeed
            .get_order(order_id)
            .await
            .context("unable to get lightspeed order for draft")?;

        if data.order.status.is_cancelled() {
            info!(
                { lightspeed_id = data.order.id },
                "order has been cancelled"
            );
            database::set_cancelled(&state.pool, &draft_id).await?;
            return Ok(OrderStatus::Cancelled);
        }

        if !data.order.status.is_shipped() {
            warn!(
                { status = ?&data.order.status },
                "label for order with unexpected shipment status"
            );
            return Ok(data.order.status);
        }

        database::set_shipment_id(
            &state.pool,
            &draft_id,
            label
                .shipment_id
                .as_ref()
                .map(ToUuid::to_uuid)
                .transpose()?,
        )
        .await?;

        state
            .lightspeed
            .update_order_status(
                order_id,
                OrderStatus::CompletedShipped,
                ShipmentStatus::Shipped,
            )
            .await
            .context("unable to update lightspeed order status")?;

        database::set_processed(&state.pool, &draft_id).await?;

        info!(
            { lightspeed_id = data.order.id, draft_id = ?&draft_id, shipment_id = &label.shipment_id },
            "order successfully processed"
        );

        return Ok(OrderStatus::CompletedShipped);
    }

    Ok(OrderStatus::ProcessingAwaitingPayment)
}

pub async fn run_once(state: AdapterState) -> Result<()> {
    let amount = database::unprocessed_count(&state.pool).await?;
    if amount == 0 {
        info!("nothing to do.");
        return Ok(());
    }

    info!("polling {amount} unprocessed orders");

    let mut orders = database::get_unprocessed_stream(&state.pool).await?;

    while let Some(query) = orders.next().await {
        match query {
            Ok(order) => match reconcile_order_status(&state, &order).await {
                Ok(status) => debug!(
                    { lightspeed_id = &order.lightspeed_order_id, dhl_id = ?&order.dhl_draft_id, status = ?&status },
                    "done checking order status"
                ),
                Err(err) => error!(
                    { err = ?&err, order_id = &order.lightspeed_order_id  },
                    "failed to check order status"
                ),
            },
            Err(err) => {
                error!(
                    { error = &err as &dyn std::error::Error },
                    "unable to get next order from database"
                );
                continue;
            }
        };
    }

    Ok(())
}
