-- +migrate Up
-- +migrate StatementBegin

DROP TABLE IF EXISTS users, tables, reservations;
DROP TYPE IF EXISTS user_role, table_status;


-- Create Table 
CREATE TYPE user_role AS ENUM ('admin', 'customer');
CREATE TYPE table_status AS ENUM ('booked', 'available');


CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- table tables



CREATE TABLE IF NOT EXISTS tables (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(255) NOT NULL,
    capacity INT NOT NULL,
    status table_status NOT NULL DEFAULT 'available',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- table reservations

CREATE TABLE IF NOT EXISTS reservations (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    table_id INT REFERENCES tables(id) ON DELETE CASCADE,
    number_of_people INT NOT NULL,
    date_reservation TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate StatementEnd