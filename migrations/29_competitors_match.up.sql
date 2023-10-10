CREATE TABLE IF NOT EXISTS competitor_match (
    match_id UUID REFERENCES matches (id) MATCH FULL ON DELETE CASCADE ON UPDATE NO ACTION DEFERRABLE,
    competitor_id UUID REFERENCES competitors (id) MATCH FULL ON DELETE CASCADE ON UPDATE NO ACTION DEFERRABLE,
    UNIQUE (match_id, competitor_id)
);
