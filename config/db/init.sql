CREATE DATABASE ingestor_db;

CONNECT ingestor_db;

CREATE TABLE customers (
    id serial PRIMARY KEY,
    name varchar,
    email varchar UNIQUE,
    mobile_number varchar,
    created_at timestamptz DEFAULT current_timestamp
  )
