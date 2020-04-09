CREATE TABLE IF NOT EXISTS links(
       id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
       link VARCHAR(1024) NOT NULL,
       icon_path VARCHAR(2048) NOT NULL,
       display_name VARCHAR(128),
       created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);