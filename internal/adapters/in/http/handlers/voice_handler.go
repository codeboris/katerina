package handlers

import (
	"io"
	"net/http"

	"github.com/codeboris/katerina/internal/adapters/in/http/middleware"
	"github.com/codeboris/katerina/internal/core/application/usecases"
)

type VoiceHandler struct {
	processUC *usecases.ProcessVoiceUseCase
}

func NewVoiceHandler(processUC *usecases.ProcessVoiceUseCase) *VoiceHandler {
	return &VoiceHandler{processUC: processUC}
}

func (h *VoiceHandler) Process(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid form data")
		return
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		writeError(w, http.StatusBadRequest, "audio file required")
		return
	}
	defer file.Close()

	audioData, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read audio failed")
		return
	}

	out, err := h.processUC.Execute(r.Context(), usecases.ProcessVoiceInput{
		UserID:    userID,
		AudioData: audioData,
		MimeType:  header.Header.Get("Content-Type"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "processing failed: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, out)
}
