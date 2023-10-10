package repository

import (
	"github.com/gofrs/uuid"
)

// NewSQLCompatUUIDFromStr Returns SQL compatible uuid for less boilerplate from string.
func NewSQLCompatUUIDFromStr(id string) uuid.NullUUID {
	return uuid.NullUUID{
		Valid: true,
		UUID:  uuid.FromStringOrNil(id),
	}
}

func NewSQLCompatUUIDNULL() uuid.NullUUID {
	return uuid.NullUUID{
		Valid: false,
		UUID:  uuid.Nil,
	}
}
