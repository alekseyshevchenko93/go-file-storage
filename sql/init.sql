CREATE TABLE IF NOT EXISTS files (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  key TEXT UNIQUE,
  extension TEXT,
  last_downloaded_at TIMESTAMPTZ DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS files_key_idx ON files USING btree (key);
CREATE INDEX IF NOT EXISTS files_last_downloaded_at_idx ON files USING btree (last_downloaded_at);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_files_updated_at
  BEFORE UPDATE
  ON
    files
  FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
