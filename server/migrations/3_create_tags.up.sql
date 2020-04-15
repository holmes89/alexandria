CREATE TABLE IF NOT EXISTS tags(
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  display_name VARCHAR(128) UNIQUE NOT NULL,
  color VARCHAR(8)
);

CREATE TABLE IF NOT EXISTS tagged_resources(
  id uuid NOT NULL,
  resource_id uuid NOT NULL,
  resource_type VARCHAR(32),
  FOREIGN KEY (id) REFERENCES tags (id)
);