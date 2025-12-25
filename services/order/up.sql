CREATE TABLE IF NOT EXISTS orders (
  id SERIAL PRIMARY KEY,
  account_id INTEGER NOT NULL,
  total_price BIGINT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_products (
  order_id INTEGER REFERENCES orders (id) ON DELETE CASCADE,
  product_id TEXT NOT NULL,
  quantity INT NOT NULL,
  PRIMARY KEY (order_id, product_id)
);