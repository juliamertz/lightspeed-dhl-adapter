use anyhow::{Context, Result};
use futures::StreamExt;
use tracing::{debug, error, info};

use crate::{
    AdapterState,
    database::{self, models::Order},
    lightspeed::{OrderStatus, ShipmentStatus},
    utils::ToUuid,
};

async fn reconcile_order_status(state: &AdapterState, order: &Order) -> Result<OrderStatus> {
    let draft_id = order
        .dhl_draft_id
        .context("order has been created without a draft id")?;

    database::incr_poll_count(&state.pool, order.lightspeed_order_id).await?;

    debug!(
        { dhl_id = ?&order.dhl_draft_id, lightspeed_id = &order.lightspeed_order_id },
        "checking order's shipment status"
    );

    // if a label has been created for this order 
    // we can go ahead and update the status in lightspeed
    if let Some(label) = state
        .dhl
        .get_label(&order.lightspeed_order_id.to_string())
        .await?
    {
        info!({ shipment_id = label.shipment_id }, "found label for order");

        let shipment_id = label
            .shipment_id
            .as_ref()
            .context("label has no shipment id set")?
            .to_uuid()?;

        database::set_shipment_id(&state.pool, &draft_id, &shipment_id).await?;

        let data = state
            .lightspeed
            .get_order(&order.lightspeed_order_id.to_string())
            .await
            .context("unable to get lightspeed order for draft")?;

        if let OrderStatus::Cancelled = data.order.status {
            info!(
                { lightspeed_id = data.order.id },
                "order has been cancelled"
            );
            database::set_cancelled(&state.pool, &draft_id).await?;
            return Ok(OrderStatus::Cancelled);
        }

        state
            .lightspeed
            .update_order_status(
                &order.lightspeed_order_id.to_string(),
                OrderStatus::CompletedShipped,
                ShipmentStatus::Shipped,
            )
            .await
            .context("unable to update lightspeed order status")?;

        database::set_processed(&state.pool, &draft_id).await?;
        return Ok(OrderStatus::CompletedShipped)
    }

    Ok(OrderStatus::ProcessingAwaitingPayment)
}

pub async fn run_once(state: AdapterState) -> Result<()> {
    debug!("polling unprocessed orders");

    let unprocessed_count = database::unprocessed_count(&state.pool).await?;
    info!("found {unprocessed_count} unprocessed orders");

    if unprocessed_count == 0 {
        info!("nothing to do.");
        return Ok(());
    }

    let mut orders = database::get_unprocessed_stream(&state.pool).await?;

    while let Some(query) = orders.next().await {
        match query {
            Ok(order) => match reconcile_order_status(&state, &order).await {
                Ok(status) => info!(
                    { order_id = &order.lightspeed_order_id, status = ?&status },
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
