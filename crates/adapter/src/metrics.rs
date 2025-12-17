use std::net::SocketAddr;

use axum::Router;
use axum::extract::State;
use axum::response::{IntoResponse, Response};
use axum::routing::get;
use axum_macros::FromRef;
use axum_prometheus::PrometheusMetricLayer;
use axum_prometheus::metrics_exporter_prometheus::PrometheusHandle;
use metrics::gauge;
use tokio::net::TcpListener;

use crate::database::{self, ConnectionPool};
use crate::{AdapterError, AdapterState};

#[derive(Clone, FromRef)]
struct MetricState {
    handle: PrometheusHandle,
    pool: ConnectionPool,
}

pub async fn serve(addr: SocketAddr, state: AdapterState) {
    let (_, handle) = PrometheusMetricLayer::pair();
    let state = MetricState {
        handle,
        pool: state.pool,
    };

    let app = Router::new()
        .route("/metrics", get(metrics))
        .with_state(state);

    let listener = TcpListener::bind(&addr).await.unwrap();
    tracing::info!("listening on {}", listener.local_addr().unwrap());
    axum::serve(listener, app).await.unwrap();
}

async fn metrics(
    State(handle): State<PrometheusHandle>,
    State(pool): State<ConnectionPool>,
) -> Result<Response, AdapterError> {
    let conn = &mut pool.get().await?;

    gauge!("lightspeed_dhl_adapter_unprocessed_count")
        .set(database::unprocessed_count(conn).await? as f64);
    gauge!("lightspeed_dhl_adapter_processed_count")
        .set(database::processed_count(conn).await? as f64);

    Ok(handle.render().into_response())
}
