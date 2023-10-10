package models

import (
	"github.com/gofrs/uuid"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

// Competitor model
type CompetitorMatch struct {
	TableName    struct{}      `sql:"competitor_match"`
	MatchID      uuid.NullUUID `pg:",fk",discard_unknown_columns`
	CompetitorID uuid.NullUUID `pg:",fk",discard_unknown_columns`
	Competitor   *Competitor
	Match        *Match
}

func (t *CompetitorMatch) GetMatchID() string {
	return t.MatchID.UUID.String()
}

func (t *CompetitorMatch) GetCompetitorID() string {
	return t.CompetitorID.UUID.String()
}

// Tests Equality of parts of the model that are not child objects or the ID
func (t *CompetitorMatch) EqualsShallow(b *CompetitorMatch) bool {
	return (t.MatchID == b.MatchID &&
		t.CompetitorID == b.CompetitorID)
}

func (t *CompetitorMatch) AddRelationship(c *Competitor, m *Match, db repository.DataSource) error {

	competitorMatch := CompetitorMatch{}
	count, err := db.Model(&competitorMatch).Where("match_id = ? and competitor_id = ?", m.ID, c.ID).SelectAndCount()

	if count == 0 {
		competitorMatch = CompetitorMatch{MatchID: m.ID, CompetitorID: c.ID}
		_, err = db.Model(&competitorMatch).Returning("*").Insert()
		if err != nil {
			return err
		}
	}

	return err
}
