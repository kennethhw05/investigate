ALTER TABLE pools
    DROP COLUMN IF EXISTS unit_value,
    DROP COLUMN IF EXISTS min_unit_per_line,
    DROP COLUMN IF EXISTS max_unit_per_line,
    DROP COLUMN IF EXISTS min_unit_per_ticket,
    DROP COLUMN IF EXISTS max_unit_per_ticket,
    DROP COLUMN IF EXISTS currency;
