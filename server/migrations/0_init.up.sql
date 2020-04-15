CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS documents(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    description VARCHAR(1024),
    displayname VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    created timestamp NOT NULL DEFAULT current_timestamp,
    updated timestamp NULL DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    username VARCHAR(128) NOT NULL,
    password VARCHAR(1024) NOT NULL
);