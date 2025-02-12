CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(256) NOT NULL,
    balance INT DEFAULT 1000 NOT NULL
);

CREATE TABLE items(
    id SERIAL PRIMARY KEY,
    name VARCHAR(16) NOT NULL UNIQUE,
    price INT NOT NULL
);

CREATE TABLE sales(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    item_id INT NOT NULL REFERENCES items(id),
    quantity INT NOT NULL DEFAULT 1,
    UNIQUE (user_id, item_id)
);

CREATE TABLE operations(
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL REFERENCES users(id),
    receiver_id INT NOT NULL REFERENCES users(id),
    amount INT NOT NULL,
    UNIQUE (sender_id, receiver_id)
);

-- TODO: Индексы (сначала составить запросы)

INSERT INTO items(name, price) VALUES
                                   ('t-shirt', 80),
                                   ('cup', 20),
                                   ('book', 50),
                                   ('pen', 10),
                                   ('powerbank', 200),
                                   ('hoody', 300),
                                   ('umbrella', 200),
                                   ('socks', 10),
                                   ('wallet', 50),
                                   ('pink-hoody', 500);