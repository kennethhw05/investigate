package rest

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/auth"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

// SystemStateHandler handles restful calls for system state
type SystemStateHandler struct {
	CFG    *config.Config
	DB     repository.DataSource
	Logger *logrus.Logger
}

type pumpActive struct {
	IsActive bool `json:"isFeedActive"`
}

// UpdatePumpActiveHandler changes the system feed isActive status
func (h *SystemStateHandler) UpdatePumpActiveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var requestedState pumpActive
	err = decoder.Decode(&requestedState)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var isActive bool

	_, err = h.DB.QueryOne(&isActive, "update system_state set is_feed_active = ? returning is_feed_active", requestedState.IsActive)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pumpActive{
			IsActive: isActive,
		})
	}
}

// GetPumpActiveHandler responds with current system feed isActive status
func (h *SystemStateHandler) GetPumpActiveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
	})

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var isActive bool
	_, err = h.DB.QueryOne(&isActive, "SELECT system_state.is_feed_active from system_state limit 1;")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(pumpActive{
			IsActive: isActive,
		})
	}
}
