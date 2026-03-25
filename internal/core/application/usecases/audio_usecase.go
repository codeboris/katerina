package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/codeboris/katerina/internal/core/ports"
)

var ErrAudioNotFound = errors.New("audio not found")

type AudioUseCase struct {
	audioRepo ports.AudioRepository
}

func NewAudioUseCase(audioRepo ports.AudioRepository) *AudioUseCase {
	return &AudioUseCase{audioRepo: audioRepo}
}

func (uc *AudioUseCase) GetAudioPath(ctx context.Context, id uuid.UUID) (string, error) {
	audio, err := uc.audioRepo.FindByID(ctx, id)
	if err != nil {
		return "", ErrAudioNotFound
	}
	return audio.FilePath, nil
}
