use base64::{Engine as _, prelude::BASE64_STANDARD as base64};
use reqwest::{Client as HttpClient, header::HeaderMap};
use serde::Serialize;
use thiserror::Error;
use tracing::instrument;

use crate::models::{CatalogResponse, Options, OrderStatus, OrderWrapper, ShipmentStatus};

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
}

#[derive(Serialize)]
#[serde(rename_all = "camelCase")]
struct UpdateStatusRequest<'a> {
    status: &'a OrderStatus,
    shipment_status: &'a ShipmentStatus,
}

impl LightspeedClient {
    pub fn new(options: Options) -> Self {
        Self {
            options,
            http: reqwest::Client::new(),
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

    #[instrument(skip(self), err(Debug))]
    pub async fn get_order(&self, id: u64) -> Result<OrderWrapper> {
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

    #[instrument(skip(self), err(Debug))]
    pub async fn update_order_status(
        &self,
        id: u64,
        status: OrderStatus,
        shipment_status: ShipmentStatus,
    ) -> Result<OrderWrapper> {
        let request = UpdateStatusRequest {
            status: &status,
            shipment_status: &shipment_status,
        };

        Ok(self
            .http
            .put(self.endpoint(format!("orders/{id}.json")))
            .json(&OrderWrapper::new(request))
            .headers(self.headers()?)
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    #[instrument(skip(self), err(Debug))]
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
