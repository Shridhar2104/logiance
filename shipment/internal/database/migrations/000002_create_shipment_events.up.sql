-- migrations/000002_create_shipment_events.up.sql
CREATE TABLE IF NOT EXISTS shipment_events (
    id BIGSERIAL PRIMARY KEY,
    shipment_tracking_id BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL,
    location VARCHAR(255),
    timestamp TIMESTAMP NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_shipment_tracking
        FOREIGN KEY (shipment_tracking_id)
        REFERENCES shipment_trackings(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_shipment_events_tracking_id ON shipment_events(shipment_tracking_id);
CREATE INDEX idx_shipment_events_status ON shipment_events(status);
CREATE INDEX idx_shipment_events_timestamp ON shipment_events(timestamp);
