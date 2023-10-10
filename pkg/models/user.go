package models

import (
	"github.com/gofrs/uuid"
)

type User struct {
	TableName  struct{} `sql:"users,alias:user"`
	Password   []byte
	Email      string
	Token      string
	ID         uuid.NullUUID
	AccessRole AccessRole
	ResetToken string
}

func (u *User) GetID() string {
	return u.ID.UUID.String()
}
