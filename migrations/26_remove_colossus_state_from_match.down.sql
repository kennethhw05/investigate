ALTER TABLE matches ADD COLUMN IF NOT EXISTS synced_colossus_status TEXT DEFAULT 'UNKNOWN';