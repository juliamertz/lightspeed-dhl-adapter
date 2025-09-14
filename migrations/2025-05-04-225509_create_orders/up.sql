CREATE TABLE orders (
  id             SERIAL PRIMARY KEY NOT NULL,
  incoming_order JSONB,
  poll_count     INTEGER DEFAULT 0 NOT NULL,
  created_at     TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at     TIMESTAMP NOT NULL DEFAULT current_timestamp,
  cancelled_at   TIMESTAMP,
  processed_at   TIMESTAMP,

  lightspeed_order_id     INTEGER NOT NULL,
  lightspeed_order_number TEXT NOT NULL,

  dhl_draft_id    UUID NOT NULL,
  dhl_shipment_id UUID
);

CREATE INDEX idx_orders_dhl_id 
ON orders (dhl_draft_id);

CREATE INDEX idx_lightspeed_orders_number 
ON orders (lightspeed_order_number);

CREATE INDEX idx_orders_unprocessed 
ON orders (created_at)
WHERE processed_at IS NULL AND cancelled_at IS NULL;

