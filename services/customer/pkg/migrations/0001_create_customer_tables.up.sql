CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    category_id INT,
    name VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    cost DECIMAL(10, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS bucket (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    preparing BOOL NOT NULL
);

CREATE TABLE IF NOT EXISTS bucket_items (
    id SERIAL PRIMARY KEY,
    bucket_id INT NOT NULL,
    product_id INT NOT NULL,
    amount INT NOT NULL,
    UNIQUE (bucket_id, product_id)
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    bucket_id INT NOT NULL,
    total_cost DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delivered_at TIMESTAMP,
    status VARCHAR(255) NOT NULL
);