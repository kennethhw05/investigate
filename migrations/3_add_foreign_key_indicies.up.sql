CREATE INDEX IF NOT EXISTS matches_event_id ON matches USING btree (event_id);
CREATE INDEX IF NOT EXISTS teams_match_id ON teams USING btree (match_id);
CREATE INDEX IF NOT EXISTS players_team_id ON players USING btree (team_id);
CREATE INDEX IF NOT EXISTS legs_match_id ON legs USING btree (match_id);
CREATE INDEX IF NOT EXISTS legs_pool_id ON legs USING btree (pool_id);
