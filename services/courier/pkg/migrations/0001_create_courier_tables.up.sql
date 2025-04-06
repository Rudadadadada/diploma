CREATE TABLE IF NOT EXISTS couriers (
    id INT PRIMARY KEY,
    active BOOL NOT NULL,
    in_progress BOOL NOT NULL,
    rating DECIMAL(10, 2) NOT NULL,
    order_delivered INT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id INT PRIMARY KEY,
    customer_id INT NOT NULL,
    courier_id INT NOT NULL,
    total_cost DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP,
    delivery_started TIMESTAMP,
    delivered_at TIMESTAMP,
    took BOOL NOT NULL,
    status VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS order_items (
    id INT PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    total_cost DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);