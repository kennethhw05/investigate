package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/auth"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

// MatchHandler handles restful calls for matches
type MatchHandler struct {
	CFG    *config.Config
	DB     repository.DataSource
	Logger *logrus.Logger
}

// MatchStatisticsHandler responds with match statistics json blob
func (h *MatchHandler) MatchStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
		models.AccessRoleGuestAPIOnly,
	})

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	matchID := chi.URLParam(r, "matchID")
	match := models.Match{}

	err = h.DB.Model(&match).Where("id = ?", repository.NewSQLCompatUUIDFromStr(matchID)).Select()

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"match_%s_statistics.json\"", matchID))
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	json.NewEncoder(w).Encode(match.Statistics)
}
