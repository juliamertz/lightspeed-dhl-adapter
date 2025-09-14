use axum::http::HeaderMap;
use base64::{Engine as _, prelude::BASE64_STANDARD as base64};
use reqwest::Client as HttpClient;
use serde_json::json;
use thiserror::Error;
use tracing::info;

use crate::lightspeed::models::{
    CatalogResponse, IncomingOrder, Options, OrderStatus, ShipmentStatus,
};

use super::models::Product;

#[derive(Debug, Error)]
pub enum LightspeedError {
    #[error("invalid json: '{0}'")]
    Json(#[from] serde_json::Error),
    #[error("http request failed: '{0}'")]
    Http(#[from] reqwest::Error),
    #[error("invalid header value: '{0}'")]
    InvalidHeaderValue(#[from] reqwest::header::InvalidHeaderValue),
}

type Result<T> = std::result::Result<T, LightspeedError>;

#[derive(Debug, Clone)]
pub struct LightspeedClient {
    options: Options,
    http: HttpClient,
    skip_mutation: bool,
}

impl LightspeedClient {
    pub fn new(options: Options, skip_mutation: bool) -> Self {
        Self {
            options,
            http: reqwest::Client::new(),
            skip_mutation,
        }
    }

    fn endpoint(&self, path: impl AsRef<str>) -> String {
        format!(
            "{url}/{path}",
            url = self.options.cluster,
            path = path.as_ref()
        )
    }

    fn headers(&self) -> Result<HeaderMap> {
        let credentials = format!(
            "{username}:{password}",
            username = self.options.key,
            password = self.options.secret
        );
        let basic_auth = format!("Basic {}", base64.encode(credentials));

        let mut headers = HeaderMap::new();
        headers.insert("Authorization", basic_auth.try_into()?);
        headers.insert(
            "Content-Type",
            "application/json".try_into().expect("valid header"),
        );
        headers.insert(
            "Accept",
            "application/json".try_into().expect("valid header"),
        );
        Ok(headers)
    }

    pub async fn get_order(&self, id: &str) -> Result<IncomingOrder> {
        Ok(self
            .http
            .get(self.endpoint(format!("orders/{id}.json")))
            .headers(self.headers()?)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    pub async fn update_order_status(
        &self,
        id: &str,
        status: OrderStatus,
        shipment_status: ShipmentStatus,
    ) -> Result<()> {
        if self.skip_mutation {
            info!({ order_id = id, status = ?&status }, "mocking order status update");
            return Ok(());
        }

        let body = json!({
            "order": {
                "status": serde_json::to_string(&status)?,
                "shipmentStatus": serde_json::to_string(&shipment_status)?,
            }
        });

        Ok(self
            .http
            .put(self.endpoint(format!("orders/{id}.json")))
            .json(&body)
            .headers(self.headers()?)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    pub async fn get_stock_under_threshold(&self) -> Result<Vec<Product>> {
        let response: CatalogResponse = self
            .http
            .get(self.endpoint("catalog.json"))
            .headers(self.headers()?)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;

        Ok(response
            .products
            .into_iter()
            .filter(|product| {
                product
                    .variants
                    .iter()
                    .any(|(_, variant)| variant.stock_level <= variant.stock_alert)
            })
            .collect())
    }
}
