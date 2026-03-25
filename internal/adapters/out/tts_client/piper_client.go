package tts_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type PiperClient struct {
	baseURL      string
	audioStorage string
	client       *http.Client
}

func NewPiperClient(baseURL, audioStorage string) *PiperClient {
	return &PiperClient{
		baseURL:      baseURL,
		audioStorage: audioStorage,
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

type ttsRequest struct {
	Text  string `json:"text"`
	Voice string `json:"voice"`
}

func (c *PiperClient) Synthesize(ctx context.Context, text, voice string) (string, error) {
	body, _ := json.Marshal(ttsRequest{Text: text, Voice: voice})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/synthesize", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("tts request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("tts returned status %d", resp.StatusCode)
	}

	if err := os.MkdirAll(c.audioStorage, 0755); err != nil {
		return "", fmt.Errorf("create audio dir: %w", err)
	}

	filename := fmt.Sprintf("%s.wav", uuid.New().String())
	filePath := filepath.Join(c.audioStorage, filename)

	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("create audio file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", fmt.Errorf("write audio: %w", err)
	}
	return filePath, nil
}
