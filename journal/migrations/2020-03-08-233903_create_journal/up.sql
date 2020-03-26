CREATE EXTENSION IF NOT EXISTS pgcrypto;
DROP TABLE IF EXISTS journal_entry;

CREATE TABLE journal_entry (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  content TEXT NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
