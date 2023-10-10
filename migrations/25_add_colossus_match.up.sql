CREATE TABLE IF NOT EXISTS colossus_matches (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_id TEXT,
    status TEXT,
    pool_type TEXT,
    match_id UUID REFERENCES matches (id) MATCH FULL ON DELETE CASCADE ON UPDATE NO ACTION DEFERRABLE,
    UNIQUE (match_id, pool_type)
);
