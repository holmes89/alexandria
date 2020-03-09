CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE journal_entry {
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  content TEXT NOT NULL,
  created DEFAULT TIMESTAMP NOT NULL
}
