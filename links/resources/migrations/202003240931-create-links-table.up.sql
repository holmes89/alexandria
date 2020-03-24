CREATE EXTENSION pgcrypto;
--;;
CREATE TABLE IF NOT EXISTS links(
       id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
       link VARCHAR(1024) NOT NULL,
       display_name VARCHAR(128),
       description TEXT,
       created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
