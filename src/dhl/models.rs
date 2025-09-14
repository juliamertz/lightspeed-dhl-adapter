use serde::{Deserialize, Serialize};
use serde_with::skip_serializing_none;

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all(serialize = "camelCase"))]
pub struct Credentials {
    pub user_id: String,
    #[serde(rename(serialize = "key"))]
    pub api_key: String,
    #[allow(unused)]
    #[serde(skip_serializing)]
    pub account_id: String,
}

#[derive(Debug, Deserialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct ApiToken {
    pub access_token: String,
    pub access_token_expiration: i64,
    pub refresh_token: String,
    pub refresh_token_expiration: i64,
    pub account_numbers: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct Draft {
    pub id: String,
    pub shipment_id: String,
    pub order_reference: String,
    pub receiver: Option<Contact>,
    pub bcc_emails: Option<Vec<String>>,
    pub shipper: Option<Shipper>,
    pub account_id: Option<String>,
    pub options: Vec<OptionField>,
    pub on_behalf_of: Option<Shipper>,
    pub product: Option<String>,
    pub customs_declaration: Option<CustomsDeclaration>,
    pub return_label: bool,
    pub pieces: Vec<Piece>,
    pub metadata: Option<Metadata>,
    // delivery_area: Option<DeliveryArea>,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct Contact {
    pub name: Option<Name>,
    pub address: Option<Address>,
    pub email: Option<String>,
    pub phone_number: Option<String>,
    pub vat_number: Option<String>,
    pub eori_number: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct Shipper {
    pub name: Option<Name>,
    pub address: Option<Address>,
    pub email: Option<String>,
    pub phone_number: Option<String>,
    pub vat_number: Option<String>,
    pub eori_number: Option<String>,
    pub rex_number: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct Name {
    pub first_name: Option<String>,
    pub last_name: Option<String>,
    pub company_name: Option<String>,
    pub additional_name: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct Address {
    pub country_code: Option<String>,
    pub postal_code: Option<String>,
    pub city: Option<String>,
    pub street: Option<String>,
    pub additional_address_line: Option<String>,
    pub number: Option<String>,
    pub is_business: bool,
    pub addition: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct OptionField {
    pub key: String,
    pub input: String,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct Label {
    pub label_id: Option<String>,
    pub order_reference: Option<String>,
    pub parcel_type: Option<String>,
    pub label_type: Option<String>,
    pub piece_number: Option<i32>,
    pub tracker_code: Option<String>,
    pub routing_code: Option<String>,
    pub user_id: Option<String>,
    pub organisation_id: Option<String>,
    pub application: Option<String>,
    pub time_created: Option<String>,
    pub shipment_id: Option<String>,
    pub account_number: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct CustomsDeclaration {
    pub certificate_number: Option<String>,
    pub currency: Option<String>,
    pub invoice_number: Option<String>,
    pub licence_number: Option<String>,
    pub remarks: Option<String>,
    pub invoice_type: Option<String>,
    pub export_type: Option<String>,
    pub export_reason: Option<String>,
    pub customs_goods: Option<Vec<CustomsGood>>,
    pub inco_terms: Option<String>,
    pub inco_terms_city: Option<String>,
    pub sender_inbound_vat_number: Option<String>,
    pub attachment_ids: Option<Vec<String>>,
    pub shipping_fee: Option<ShippingFee>,
    pub importer_of_record: Option<Importer>,
    pub deferment_account_vat: Option<String>,
    pub deferment_account_duties: Option<String>,
    pub vat_reverse_charge: bool,
    pub sender_has_inbound_eori_number: bool,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct CustomsGood {
    pub code: Option<String>,
    pub description: Option<String>,
    pub origin: Option<String>,
    pub quantity: Option<i32>,
    pub value: Option<f64>,
    pub weight: Option<f64>,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct ShippingFee {
    pub value: Option<f64>,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct Importer {
    pub name: Name,
    pub address: Option<Address>,
    pub email: Option<String>,
    pub phone_number: Option<String>,
    pub vat_number: Option<String>,
    pub eori_number: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "UPPERCASE")]
pub enum ParcelType {
    Small,
    Medium,
    Large,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct Piece {
    pub parcel_type: Option<ParcelType>,
    pub quantity: Option<i32>,
    pub weight: Option<i64>,
    pub dimensions: Option<Dimensions>,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct Dimensions {
    pub length: Option<f64>,
    pub width: Option<f64>,
    pub height: Option<f64>,
}

#[allow(unused)]
#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct DeliveryArea {
    #[serde(rename = "type")]
    pub ty: Option<String>,
    pub remote: bool,
}

#[derive(Debug, Serialize, Deserialize)]
#[skip_serializing_none]
#[serde(rename_all = "camelCase")]
pub struct Metadata {
    pub import_id: Option<String>,
    pub user_id: Option<String>,
    pub last_modified_by: Option<String>,
    pub organization_id: Option<String>,
    pub time_created: Option<String>,
}
