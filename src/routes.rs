use std::{net::SocketAddr, str::FromStr, sync::Arc};

use axum::{
    Json, Router,
    extract::State,
    http::{HeaderMap, StatusCode},
    response::{IntoResponse, Response},
    routing::{get, post},
};
use serde_json::json;
use tokio::net::TcpListener;
use tower_http::trace::TraceLayer;
use tracing::instrument;
use uuid::Uuid;

use crate::AdapterError;
use crate::database::{self, ConnectionPool};
use crate::dhl::client::DHLClient;
use crate::lightspeed::{IncomingOrder, LightspeedClient};
use crate::{AdapterState, config::Config, dhl::Draft};

pub async fn serve(addr: SocketAddr, state: AdapterState) {
    let listener = TcpListener::bind(&addr).await.unwrap();
    let app = Router::new()
        .route("/ready", get(ready))
        .route("/webhook", post(webhook))
        .route("/stock-under-threshold", get(stock_under_threshold))
        .layer(TraceLayer::new_for_http())
        .with_state(state);

    tracing::info!("listening on {}", listener.local_addr().unwrap());
    axum::serve(listener, app).await.unwrap();
}

pub async fn ready() -> Response {
    let body = json!({ "ready": true });
    Json(body).into_response()
}

#[instrument(err(Debug))]
pub async fn webhook(
    headers: HeaderMap,
    State(pool): State<ConnectionPool>,
    State(config): State<Arc<Config>>,
    State(dhl): State<Arc<DHLClient>>,
    Json(incoming): Json<IncomingOrder>,
) -> Result<Response, AdapterError> {
    let (Some(cluster_id), Some(shop_id)) = (headers.get("x-cluster-id"), headers.get("x-shop-id"))
    else {
        return Ok(StatusCode::BAD_REQUEST.into_response());
    };

    if cluster_id.to_str()? != config.lightspeed.cluster_id
        || shop_id.to_str()? != config.lightspeed.shop_id
    {
        tracing::warn!("attempt to call '/webhook' with invalid authorization");
        return Ok(StatusCode::UNAUTHORIZED.into_response());
    }

    tracing::debug!(
        { incoming = serde_json::to_string(&incoming)? },
        "received webhook"
    );

    database::create_order(&pool, &incoming).await?;

    let draft: Draft = incoming.clone().into();
    dhl.create_draft(&draft).await?;

    let draft_id = Uuid::from_str(&draft.id)?;
    database::link_dhl_draft(&pool, &draft_id, incoming.order.id as i32).await?;

    Ok(StatusCode::OK.into_response())
}

#[instrument(err(Debug))]
pub async fn stock_under_threshold(
    State(lightspeed): State<Arc<LightspeedClient>>,
) -> Result<Response, AdapterError> {
    let stock = lightspeed.get_stock_under_threshold().await?;
    Ok(Json(stock).into_response())
}
