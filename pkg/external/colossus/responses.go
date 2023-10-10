package colossus

import (
	"time"
)

type GenericResponse struct {
	Message string
}

type PoolStatusResponse struct {
	Status           *PoolStatus `json:"status"`
	SettlementStatus *string     `json:"settlement_status"`
	SettledAt        *time.Time  `json:"settled_at"`
}

type SportEventStatusResponse struct {
	Status *SportEventStatus `json:"status"`
	Period *int              `json:"period"`
}
