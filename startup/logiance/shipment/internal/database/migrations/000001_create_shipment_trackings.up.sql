CREATE TABLE IF NOT EXISTS shipment_trackings (
    id BIGSERIAL PRIMARY KEY,
    account_id VARCHAR(255) NOT NULL,
    order_id VARCHAR(255) UNIQUE NOT NULL,
    tracking_id VARCHAR(255) UNIQUE NOT NULL,
    awb_number VARCHAR(255) UNIQUE NOT NULL,
    courier_code VARCHAR(50) NOT NULL,
    status VARCHAR(50),
    label TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);
CREATE INDEX idx_shipment_trackings_account_id ON shipment_trackings(account_id);
CREATE INDEX idx_shipment_trackings_order_id ON shipment_trackings(order_id);
CREATE INDEX idx_shipment_trackings_tracking_id ON shipment_trackings(tracking_id);
CREATE INDEX idx_shipment_trackings_awb ON shipment_trackings(awb_number);
CREATE INDEX idx_shipment_trackings_courier ON shipment_trackings(courier_code);
CREATE INDEX idx_shipment_trackings_status ON shipment_trackings(status);
