package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Audit struct {
	ID         uuid.NullUUID
	Time       time.Time
	TargetID   uuid.NullUUID
	TargetType string
	UserID     uuid.NullUUID `pg:",fk"`
	Content    string
	EditAction EditAction
}

func (audit *Audit) GetID() string {
	return audit.ID.UUID.String()
}

func (audit *Audit) GetTargetID() string {
	return audit.TargetID.UUID.String()
}

func (audit *Audit) GetUserID() string {
	return audit.UserID.UUID.String()
}
