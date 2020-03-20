CREATE TABLE IF NOT EXISTS documents(
    id VARCHAR(36) NOT NULL,
    description VARCHAR(1024),
    displayName VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    created timestamp NOT NULL DEFAULT current_timestamp,
    updated timestamp NULL DEFAULT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) NOT NULL,
    username VARCHAR(128) NOT NULL,
    password VARCHAR(1024) NOT NULL,
    PRIMARY KEY(id)
);