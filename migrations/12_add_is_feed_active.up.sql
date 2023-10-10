CREATE TABLE IF NOT EXISTS system_state (
    is_feed_active BOOLEAN
);

INSERT INTO system_state (is_feed_active) values (TRUE);
