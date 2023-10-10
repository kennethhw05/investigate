package models

import (
	"github.com/gofrs/uuid"
)

// Player struct for a player in our system
type Player struct {
	ExternalID string
	Name       string
	Nickname   string
	ID         uuid.NullUUID
	TeamID     uuid.NullUUID `pg:",fk"`
	TableName  struct{}      `sql:"players,alias:player"`
}

func (p *Player) GetID() string {
	return p.ID.UUID.String()
}

func (p *Player) GetTeamID() string {
	return p.TeamID.UUID.String()
}

// Tests Equality of parts of the model that are not child objects or the ID
func (p *Player) EqualsShallow(b *Player) bool {
	return (p.ExternalID == b.ExternalID &&
		p.Name == b.Name &&
		p.Nickname == b.Nickname &&
		p.TeamID == b.TeamID)
}
