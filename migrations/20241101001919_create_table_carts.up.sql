CREATE TABLE IF NOT EXISTS carts (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  user_id VARCHAR(20),
  price_total NUMERIC(10, 2),
  finish INT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);
