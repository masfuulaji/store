CREATE TABLE IF NOT EXISTS cart_items (
  id SERIAL PRIMARY KEY,
  cart_id VARCHAR(20),
  product_id VARCHAR(20),
  product_qty INT,
  price_total NUMERIC(10, 2),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);
