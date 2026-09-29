USE legacy;

CREATE TABLE inventory_items (
  sku VARCHAR(64) PRIMARY KEY,
  quantity INT NOT NULL,
  CONSTRAINT quantity_nonnegative CHECK (quantity >= 0)
);

CREATE TABLE stock_movements (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  sku VARCHAR(64) NOT NULL,
  delta INT NOT NULL,
  location_code VARCHAR(64) NOT NULL,
  actor_id VARCHAR(64) NOT NULL,
  occurred_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  CONSTRAINT fk_stock_sku FOREIGN KEY (sku) REFERENCES inventory_items(sku)
);

CREATE TABLE outbox_events (
  sequence_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  event_id CHAR(36) NOT NULL UNIQUE,
  event_type VARCHAR(128) NOT NULL,
  aggregate_id VARCHAR(64) NOT NULL,
  payload JSON NOT NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
);

DELIMITER //
CREATE TRIGGER stock_movement_outbox AFTER INSERT ON stock_movements
FOR EACH ROW
BEGIN
  UPDATE inventory_items
     SET quantity = quantity + NEW.delta
   WHERE sku = NEW.sku;

  INSERT INTO outbox_events (event_id, event_type, aggregate_id, payload)
  VALUES (
    UUID(),
    'inventory.stock-moved.v1',
    NEW.sku,
    JSON_OBJECT(
      'movement_id', NEW.id,
      'sku', NEW.sku,
      'delta', NEW.delta,
      'location_code', NEW.location_code,
      'actor_id', NEW.actor_id
    )
  );
END//
DELIMITER ;

INSERT INTO inventory_items (sku, quantity) VALUES ('DEMO-001', 100);

CREATE USER 'cdc'@'%' IDENTIFIED BY 'localdemo';
GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'cdc'@'%';
FLUSH PRIVILEGES;
