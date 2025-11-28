// @generated automatically by Diesel CLI.

diesel::table! {
    orders (id) {
        id -> Int4,
        incoming_order -> Nullable<Jsonb>,
        poll_count -> Int4,
        stale -> Bool,
        created_at -> Timestamp,
        updated_at -> Timestamp,
        cancelled_at -> Nullable<Timestamp>,
        processed_at -> Nullable<Timestamp>,
        lightspeed_order_id -> Int4,
        lightspeed_order_number -> Text,
        dhl_draft_id -> Nullable<Uuid>,
        dhl_shipment_id -> Nullable<Uuid>,
    }
}
