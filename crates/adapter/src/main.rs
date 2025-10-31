mod config;
mod database;
mod metrics;
mod poll;
mod routes;
mod schema;
mod transform;

use std::{net::SocketAddr, path::PathBuf, sync::Arc};

use axum::response::IntoResponse;
use axum_macros::FromRef;
use clap::{Parser, Subcommand};
use config::Config;
use database::ConnectionPool;
use dhl::client::DHLClient;
use diesel_async::pooled_connection::bb8;
use lightspeed::{client::LightspeedClient};
use reqwest::StatusCode;
use thiserror::Error;
use tracing::error;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[derive(Parser, Clone)]
pub struct Opts {
    #[clap(long, env = "DRY_RUN")]
    pub dry_run: Option<bool>,

    #[clap(long, env = "CONFIG_PATH")]
    pub config_path: Option<PathBuf>,

    #[clap(long, env = "DATABASE_URL")]
    pub database_url: String,

    #[command(subcommand)]
    command: Command,
}

#[derive(Subcommand, Clone)]
pub enum Command {
    /// Serve webhook endpoint
    Serve {
        /// Socket address to bind on
        #[arg(long, env = "ADDR", default_value = "0.0.0.0:8000")]
        addr: SocketAddr,
    },
    ServeMetrics {
        /// Socket address to bind on
        #[arg(long, env = "ADDR", default_value = "0.0.0.0:8010")]
        addr: SocketAddr,
    },
    /// Poll DHL for order status updates
    PollStatus,
}

#[derive(Debug, Error)]
enum AdapterError {
    #[error("database error: '{0}'")]
    Database(#[from] database::DatabaseError),
    #[error("database pool error: '{0}'")]
    DatabasePool(#[from] bb8::RunError),
    #[error("lightspeed error: '{0}'")]
    Lightspeed(#[from] lightspeed::client::LightspeedError),
    #[error("dhl error: '{0}'")]
    Dhl(#[from] dhl::client::DHLError),
    #[error("uuid error: '{0}'")]
    Uuid(#[from] uuid::Error),
    #[error("json error: '{0}'")]
    Json(#[from] serde_json::Error),
    #[error("cannot stringify header: '{0}'")]
    ToStr(#[from] reqwest::header::ToStrError),
    #[error("{0}")]
    Anyhow(#[from] anyhow::Error),
}

impl IntoResponse for AdapterError {
    fn into_response(self) -> axum::response::Response {
        error!(
            { err = &self as &dyn std::error::Error },
            "route handler error"
        );
        (StatusCode::INTERNAL_SERVER_ERROR).into_response()
    }
}

#[derive(Clone, FromRef)]
pub struct AdapterState {
    pub pool: ConnectionPool,
    pub config: Arc<Config>,
    pub options: Arc<Opts>,
    pub dhl: Arc<dhl::client::DHLClient>,
    pub lightspeed: Arc<lightspeed::client::LightspeedClient>,
}

#[tokio::main]
async fn main() -> Result<(), AdapterError> {
    let opts = Opts::parse();

    tracing_subscriber::registry()
        .with(tracing_setup::logging())
        .init();

    let pool = database::establish_connection(&opts.database_url).await?;
    let config = config::load(
        opts.config_path
            .clone()
            .unwrap_or_else(|| PathBuf::from("config.toml")),
    );

    let skip_mutation = opts.dry_run.unwrap_or_default();

    let lightspeed_client = LightspeedClient::new(config.lightspeed.clone());
    let dhl_client = DHLClient::new(config.dhl.clone(), skip_mutation);

    dhl_client.authenticate().await?;

    let state = AdapterState {
        pool,
        options: opts.clone().into(),
        config: config.to_owned().into(),
        dhl: dhl_client.into(),
        lightspeed: lightspeed_client.into(),
    };

    match opts.command {
        Command::PollStatus => poll::run_once(state).await?,
        Command::Serve { addr } => routes::serve(addr, state).await,
        Command::ServeMetrics { addr } => metrics::serve(addr, state).await,
    };

    Ok(())
}
