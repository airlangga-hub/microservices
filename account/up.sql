CREATE TABLE IF NOT EXISTS accounts(
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    hashed_password TEXT NOT NULL
);