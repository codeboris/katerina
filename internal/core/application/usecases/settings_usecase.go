package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/codeboris/katerina/internal/core/ports"
)

var availableVoices = []Voice{
	{ID: "en_US-lessac-medium", Label: "US English — Male (Lessac)"},
	{ID: "en_US-amy-medium", Label: "US English — Female (Amy)"},
	{ID: "en_GB-alan-medium", Label: "British English — Male (Alan)"},
}

type Voice struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type UserSettings struct {
	Voice string `json:"voice"`
}

type SettingsUseCase struct {
	userRepo ports.UserRepository
}

func NewSettingsUseCase(userRepo ports.UserRepository) *SettingsUseCase {
	return &SettingsUseCase{userRepo: userRepo}
}

func (uc *SettingsUseCase) AvailableVoices() []Voice {
	return availableVoices
}

func (uc *SettingsUseCase) GetSettings(ctx context.Context, userID uuid.UUID) (*UserSettings, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	voice := user.VoicePreference
	if voice == "" {
		voice = availableVoices[0].ID
	}
	return &UserSettings{Voice: voice}, nil
}

func (uc *SettingsUseCase) UpdateSettings(ctx context.Context, userID uuid.UUID, voice string) error {
	for _, v := range availableVoices {
		if v.ID == voice {
			return uc.userRepo.UpdateVoicePreference(ctx, userID, voice)
		}
	}
	return fmt.Errorf("unknown voice: %s", voice)
}
