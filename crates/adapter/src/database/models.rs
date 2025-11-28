use chrono::NaiveDateTime;
use diesel::{pg::Pg, prelude::*};
use uuid::Uuid;

#[allow(dead_code)]
#[derive(Debug, Queryable, Selectable)]
#[diesel(check_for_backend(Pg))]
#[diesel(table_name = crate::schema::orders)]
pub struct Order {
    pub id: i32,
    pub incoming_order: Option<serde_json::Value>,
    pub poll_count: i32,
    pub stale: bool,
    pub created_at: NaiveDateTime,
    pub updated_at: NaiveDateTime,
    pub processed_at: Option<NaiveDateTime>,
    pub cancelled_at: Option<NaiveDateTime>,

    pub lightspeed_order_id: i32,
    pub lightspeed_order_number: String,

    pub dhl_draft_id: Option<Uuid>,
    pub dhl_shipment_id: Option<Uuid>,
}
