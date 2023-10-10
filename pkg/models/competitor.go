package models

import (
	"github.com/gofrs/uuid"
)

// Competitor model
type Competitor struct {
	TableName  struct{} `sql:"competitors,alias:competitor"`
	ExternalID string
	Name       string
	Logo       string
	ID         uuid.NullUUID `pg:",fk"`
	Matches    []Match       `pg:"many2many:competitor_match"`
}

func (t *Competitor) GetID() string {
	return t.ID.UUID.String()
}

// Tests Equality of parts of the model that are not child objects or the ID
func (t *Competitor) EqualsShallow(b *Competitor) bool {
	return (t.ExternalID == b.ExternalID &&
		t.Name == b.Name &&
		t.Logo == b.Logo)
}
