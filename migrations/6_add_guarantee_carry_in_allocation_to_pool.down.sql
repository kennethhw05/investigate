ALTER TABLE pools
    DROP COLUMN IF EXISTS guarantee,
    DROP COLUMN IF EXISTS carry_in,
    DROP COLUMN IF EXISTS allocation;
