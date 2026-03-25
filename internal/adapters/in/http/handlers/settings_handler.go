package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/codeboris/katerina/internal/adapters/in/http/middleware"
	"github.com/codeboris/katerina/internal/core/application/usecases"
)

type SettingsHandler struct {
	settingsUC *usecases.SettingsUseCase
}

func NewSettingsHandler(settingsUC *usecases.SettingsUseCase) *SettingsHandler {
	return &SettingsHandler{settingsUC: settingsUC}
}

func (h *SettingsHandler) GetVoices(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.settingsUC.AvailableVoices())
}

func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	settings, err := h.settingsUC.GetSettings(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "get settings failed")
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (h *SettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var body struct {
		Voice string `json:"voice"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Voice == "" {
		writeError(w, http.StatusBadRequest, "voice is required")
		return
	}

	if err := h.settingsUC.UpdateSettings(r.Context(), userID, body.Voice); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
