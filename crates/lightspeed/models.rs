use clap::ValueEnum;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Options {
    pub key: String,
    pub secret: String,
    pub frontend: String,
    pub cluster: String,
    pub shop_id: String,
    pub cluster_id: String,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, ValueEnum)]
#[serde(rename_all = "snake_case")]
pub enum OrderStatus {
    ProcessingAwaitingPayment,
    ProcessingAwaitingShipment,
    OnHold,
    CompletedShipped,
    Cancelled,
}

impl OrderStatus {
    pub fn is_cancelled(&self) -> bool {
        self == &Self::Cancelled
    }

    pub fn is_shipped(&self) -> bool {
        matches!(
            self,
            Self::CompletedShipped
                | Self::ProcessingAwaitingPayment
                | Self::ProcessingAwaitingShipment
                | Self::OnHold
        )
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, ValueEnum)]
#[serde(rename_all = "snake_case")]
pub enum ShipmentStatus {
    Shipped,
    NotShipped,
    Cancelled,
}

impl ShipmentStatus {
    pub fn is_shipped(&self) -> bool {
        self == &Self::Shipped
    }
}

#[allow(unused)]
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, ValueEnum)]
#[serde(rename_all = "snake_case")]
pub enum PaymentStatus {
    Paid,
    NotPaid,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct CountryCode {
    pub code: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Order {
    pub id: u32,
    pub email: String,
    pub firstname: String,
    pub lastname: String,
    pub middlename: String,
    pub company_name: String,
    pub phone: String,
    pub shipment_title: String,
    pub number: String,
    pub is_company: bool,
    pub weight: i64,

    pub status: OrderStatus,
    pub shipment_status: ShipmentStatus,

    pub address_billing_street: String,
    pub address_billing_city: String,
    pub address_billing_zipcode: String,
    pub address_billing_country: CountryCode,
    pub address_billing_number: String,
    pub address_billing_extension: String,

    pub address_shipping_street: String,
    pub address_shipping_city: String,
    pub address_shipping_zipcode: String,
    pub address_shipping_country: CountryCode,
    pub address_shipping_number: String,
    pub address_shipping_extension: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Variant {
    pub id: i32,
    pub stock_level: i32,
    pub stock_alert: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Product {
    pub id: i32,
    pub title: String,
    pub variants: HashMap<String, Variant>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct CatalogResponse {
    pub products: Vec<Product>,
}

#[derive(Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct IncomingOrder {
    pub order: Order,
}

impl std::fmt::Debug for IncomingOrder {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.write_str(
            &serde_json::to_string(&self)
                .unwrap_or("failed to stringify incoming order".to_string()),
        )
    }
}

impl IncomingOrder {
    pub fn is_shipped(&self) -> bool {
        self.order.status.is_shipped() && self.order.shipment_status.is_shipped()
    }
}
