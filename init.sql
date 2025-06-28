-- Active: 1750844219103@@127.0.0.1@5432@postgres
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    pin VARCHAR(255) NOT NULL
);

CREATE TABLE user_balance (
    id SERIAL PRIMARY KEY,
    id_user INT REFERENCES users(id),
    balance DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TYPE type_transaction AS ENUM ('income','expense');

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    transactions_date DATE NOT NULL,
    nominal DECIMAL(10,2),
    type type_transaction NOT NULL,
    id_user INT REFERENCES users(id),
    id_other_user INT REFERENCES users(id),
    notes TEXT
);

SELECT * FROM users;
DROP TABLE users;
DROP TABLE user_balance;
DROP TABLE transactions;

UPDATE users SET name = 'davinda' WHERE id = 1;

INSERT INTO users (name, email, phone_number, password, pin)
VALUES ('abc','naldy@gmail.com','089','1234','1234');

