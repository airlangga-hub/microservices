CREATE TABLE IF NOT EXISTS accounts(
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    type TEXT DEFAULT 'buyer'
);