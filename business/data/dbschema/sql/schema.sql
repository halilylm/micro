CREATE TABLE users (
    user_id UUID,
    name TEXT,
    email TEXT UNIQUE,
    roles TEXT[],
    password_hash TEXT,
    enabled BOOLEAN,
    date_created TIMESTAMP,
    date_updated TIMESTAMP,

    PRIMARY KEY (user_id)
);

CREATE TABLE products (
    product_id UUID,
    name TEXT,
    cost INT,
    quantity INT,
    user_id UUID,
    date_created TIMESTAMP,
    date_updated TIMESTAMP,

    PRIMARY KEY (product_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE sales (
    sale_id UUID,
    user_id UUID,
    product_id UUID,
    quantity INT,
    paid INT,
    date_created TIMESTAMP,

    PRIMARY KEY (sale_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
);