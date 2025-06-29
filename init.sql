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
    transactions_date TIMESTAMP DEFAULT NOW(),
    nominal DECIMAL(10,2),
    type type_transaction NOT NULL,
    id_user INT REFERENCES users(id),
    id_other_user INT REFERENCES users(id),
    notes TEXT
);

CREATE TABLE blacklist_tokens (
    token TEXT PRIMARY KEY,
    expires_at TIMESTAMP NOT NULL
);

SELECT * FROM users;
DROP TABLE users;
DROP TABLE user_balance;
DROP TABLE transactions;

UPDATE users SET name = 'davinda' WHERE id = 1;

INSERT INTO users (name, email, phone_number, password, pin)
VALUES ('abc','naldy@gmail.com','089','1234','1234');

SELECT * FROM user_balance WHERE id_user=1 
ORDER BY created_at DESC
LIMIT 1;

SELECT * FROM users WHERE name ILIKE '%1%' OR phone_number ILIKE '%08%'
OFFSET 1
LIMIT 2;

SELECT COUNT(*) AS count FROM users WHERE name ILIKE '%%' OR phone_number ILIKE '%%';

SELECT t.transactions_date, t.nominal, t.type,  
t.notes, t.id_other_user, 
u.name AS other_user_name, 
u.email AS other_user_email, 
u.phone_number AS other_user_phone 
FROM transactions t 
JOIN users u ON u.id = t.id_other_user
WHERE t.id_user=1
ORDER BY t.transactions_date DESC;

SELECT SUM(nominal) AS total_income FROM transactions
WHERE transactions_date BETWEEN '2025-06-22' AND '2025-06-29'
AND type='income' AND id_user=1;

SELECT 
    t.transactions_date, 
    t.nominal,
    CASE 
        WHEN t.type='income' THEN 'income'
        WHEN t.type='expense' AND t.id_user=2 THEN 'expense'
        WHEN t.type='expense' AND t.id_other_user=2 THEN 'income'
    END AS type,
    t.notes,
    CASE 
        WHEN t.type='income' THEN t.id_user
        WHEN t.type='expense' AND t.id_user=2 THEN t.id_other_user
        WHEN t.type='expense' AND t.id_other_user=2 THEN t.id_user
    END AS id_other_user,
    u.name AS other_user_name,
    u.email AS other_user_email,
    u.phone_number AS other_user_phone
FROM transactions t
JOIN users u ON 
    (CASE 
        WHEN t.type='income' THEN u.id = t.id_user
        WHEN t.type='expense' AND t.id_user=2 THEN u.id = t.id_other_user
        WHEN t.type='expense' AND t.id_other_user=2 THEN u.id = t.id_user
    END)
WHERE t.id_user=2 OR t.id_other_user=2
ORDER BY t.transactions_date DESC;