package usecases

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/codeboris/katerina/internal/core/domain/model"
	"github.com/codeboris/katerina/internal/core/ports"
)

type ProcessVoiceInput struct {
	UserID    uuid.UUID
	AudioData []byte
	MimeType  string
	Voice     string
}

type ProcessVoiceOutput struct {
	Original    string `json:"original"`
	Corrected   string `json:"corrected"`
	Translation string `json:"translation"`
	Answer      string `json:"answer"`
	AudioURL    string `json:"audio_url"`
}

type ProcessVoiceUseCase struct {
	stt          ports.STTPort
	llm          ports.LLMPort
	tts          ports.TTSPort
	convRepo     ports.ConversationRepository
	audioRepo    ports.AudioRepository
	audioStorage string
}

func NewProcessVoiceUseCase(
	stt ports.STTPort,
	llm ports.LLMPort,
	tts ports.TTSPort,
	convRepo ports.ConversationRepository,
	audioRepo ports.AudioRepository,
	audioStorage string,
) *ProcessVoiceUseCase {
	return &ProcessVoiceUseCase{
		stt:          stt,
		llm:          llm,
		tts:          tts,
		convRepo:     convRepo,
		audioRepo:    audioRepo,
		audioStorage: audioStorage,
	}
}

func (uc *ProcessVoiceUseCase) Execute(ctx context.Context, input ProcessVoiceInput) (*ProcessVoiceOutput, error) {
	// 1. Write uploaded audio to a temp file
	tmpID := uuid.New().String()
	tmpInput := filepath.Join(os.TempDir(), tmpID+"_input")
	tmpWAV := filepath.Join(os.TempDir(), tmpID+".wav")

	if err := os.WriteFile(tmpInput, input.AudioData, 0600); err != nil {
		return nil, fmt.Errorf("write temp audio: %w", err)
	}
	defer os.Remove(tmpInput)
	defer os.Remove(tmpWAV)

	// 2. Convert to 16 kHz mono WAV via ffmpeg
	cmd := exec.CommandContext(ctx, "ffmpeg", "-y", "-i", tmpInput,
		"-ar", "16000", "-ac", "1", "-f", "wav", tmpWAV)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg: %w — %s", err, out)
	}

	// 3. Speech-to-text
	original, err := uc.stt.Transcribe(ctx, tmpWAV)
	if err != nil {
		return nil, fmt.Errorf("stt transcribe: %w", err)
	}

	// 4. LLM: correct grammar, translate, generate answer
	llmResp, err := uc.llm.Process(ctx, original)
	if err != nil {
		return nil, fmt.Errorf("llm process: %w", err)
	}

	// 5. Text-to-speech for the AI answer
	voice := input.Voice
	if voice == "" {
		voice = "en_US-lessac-medium"
	}
	ttsPath, err := uc.tts.Synthesize(ctx, llmResp.Answer, voice)
	if err != nil {
		return nil, fmt.Errorf("tts synthesize: %w", err)
	}

	// 6. Persist audio metadata
	audioID := uuid.New()
	if err := uc.audioRepo.Save(ctx, &model.Audio{
		ID:        audioID,
		UserID:    input.UserID,
		FilePath:  ttsPath,
		MimeType:  "audio/wav",
		CreatedAt: time.Now(),
	}); err != nil {
		return nil, fmt.Errorf("save audio: %w", err)
	}

	// 7. Persist conversation
	audioURL := fmt.Sprintf("/api/audio/%s", audioID)
	if err := uc.convRepo.Save(ctx, &model.Conversation{
		ID:          uuid.New(),
		UserID:      input.UserID,
		Original:    original,
		Corrected:   llmResp.Corrected,
		Translation: llmResp.Translation,
		Answer:      llmResp.Answer,
		AudioURL:    audioURL,
		CreatedAt:   time.Now(),
	}); err != nil {
		return nil, fmt.Errorf("save conversation: %w", err)
	}

	return &ProcessVoiceOutput{
		Original:    original,
		Corrected:   llmResp.Corrected,
		Translation: llmResp.Translation,
		Answer:      llmResp.Answer,
		AudioURL:    audioURL,
	}, nil
}
