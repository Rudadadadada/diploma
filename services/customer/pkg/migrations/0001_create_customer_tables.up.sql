CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    category_id INT,
    name VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    cost DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS bucket (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    prodcuts_with_amount INT[][],
    total_cost INT NOT NULL,
    status VARCHAR(255) NOT NULL
);