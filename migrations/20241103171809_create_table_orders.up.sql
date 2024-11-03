CREATE TABLE IF NOT EXISTS orders (
  id SERIAL PRIMARY KEY,
  cart_id VARCHAR(20),
  price_total NUMERIC(10, 2),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);
