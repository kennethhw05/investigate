CREATE TABLE IF NOT EXISTS competitors (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_id TEXT,
    name TEXT,
    logo TEXT
);
