CREATE TABLE IF NOT EXISTS couriers (
    id INT PRIMARY KEY,
    active BOOL NOT NULL,
    in_progress BOOL NOT NULL,
    rating INT NOT NULL,
    order_delivered INT NOT NULL
);