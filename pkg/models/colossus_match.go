package models

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// ColossusMatch Represents a Colossus SportEvent
type ColossusMatch struct {
	TableName  struct{}      `sql:"colossus_matches,alias:colossus_match"`
	MatchID    uuid.NullUUID `pg:",fk"`
	Match      *Match
	PoolType   PoolType
	ID         uuid.NullUUID
	Status     MatchColossusStatus
	ExternalID string
}

func (cm *ColossusMatch) GetID() string {
	return cm.ID.UUID.String()
}

func (cm *ColossusMatch) GetMatchID() string {
	return cm.MatchID.UUID.String()
}

func (cm *ColossusMatch) getPoolTypeForColossus() string {
	switch cm.PoolType {
	case PoolTypeH2h:
		return "H2H"
	case PoolTypeOverUnder:
		return "OU"
	default:
		length := len(cm.PoolType.String())
		if length >= 3 {
			return cm.PoolType.String()[0:3]
		}
		return cm.PoolType.String()[0:length]
	}
}

func (cm *ColossusMatch) GetColossusID() string {
	return fmt.Sprintf("%s-%s", cm.GetMatchID(), cm.getPoolTypeForColossus())
}

func (cm *ColossusMatch) IsNotAtEndColossusState() bool {
	return cm.Status != MatchColossusStatusOfficial && cm.Status != MatchColossusStatusAbandoned
}
