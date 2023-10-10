ALTER TABLE pools
    ADD COLUMN IF NOT EXISTS unit_value DECIMAL,
    ADD COLUMN IF NOT EXISTS min_unit_per_line DECIMAL,
    ADD COLUMN IF NOT EXISTS max_unit_per_line DECIMAL,
    ADD COLUMN IF NOT EXISTS min_unit_per_ticket DECIMAL,
    ADD COLUMN IF NOT EXISTS max_unit_per_ticket DECIMAL,
    ADD COLUMN IF NOT EXISTS currency TEXT;
