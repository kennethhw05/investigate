CREATE TABLE IF NOT EXISTS audits (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
	time      TIMESTAMP WITH TIME ZONE,
	target_id  UUID,
	target_type TEXT,
	user_id    UUID,
	content    TEXT,
	edit_action		TEXT
);


