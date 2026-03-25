package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/codeboris/katerina/internal/core/application/usecases"
)

type AudioHandler struct {
	audioUC *usecases.AudioUseCase
}

func NewAudioHandler(audioUC *usecases.AudioUseCase) *AudioHandler {
	return &AudioHandler{audioUC: audioUC}
}

func (h *AudioHandler) GetAudio(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid audio id")
		return
	}

	path, err := h.audioUC.GetAudioPath(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "audio not found")
		return
	}

	w.Header().Set("Content-Type", "audio/wav")
	http.ServeFile(w, r, path)
}
