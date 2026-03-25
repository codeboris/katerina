package llm_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/codeboris/katerina/internal/core/ports"
)

const systemPrompt = `You are a friendly English conversation teacher.
The student speaks to you in English (possibly with grammar mistakes).
Your job: correct their grammar, translate to Russian, and reply naturally to what they said — like a real conversation partner would.
The "answer" must directly respond to the content of the student's message (answer questions, react to statements, continue the topic).
Return ONLY valid JSON — no markdown, no extra text:
{
  "corrected": "<grammatically corrected version of the student's text>",
  "translation": "<Russian translation of the corrected text>",
  "answer": "<your spoken reply in English, 1-2 sentences, responding directly to what the student said>",
  "answer_translation": "<Russian translation of the answer>"
}`

type OllamaClient struct {
	baseURL string
	model   string
	client  *http.Client
}

func NewOllamaClient(baseURL, model string) *OllamaClient {
	return &OllamaClient{
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

type ollamaOptions struct {
	NumPredict int `json:"num_predict"`
}

type ollamaRequest struct {
	Model   string        `json:"model"`
	Prompt  string        `json:"prompt"`
	Stream  bool          `json:"stream"`
	Format  string        `json:"format"`
	Options ollamaOptions `json:"options"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func (c *OllamaClient) Process(ctx context.Context, text string) (*ports.LLMResponse, error) {
	prompt := systemPrompt + "\n\nUser text: " + text

	body, _ := json.Marshal(ollamaRequest{Model: c.model, Prompt: prompt, Stream: false, Format: "json", Options: ollamaOptions{NumPredict: 1024}})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ollama request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, err
	}

	jsonStr := extractJSON(ollamaResp.Response)
	var llmResp ports.LLMResponse
	if err := json.Unmarshal([]byte(jsonStr), &llmResp); err != nil {
		return nil, fmt.Errorf("parse llm JSON: %w (raw: %s)", err, ollamaResp.Response)
	}
	return &llmResp, nil
}

// extractJSON strips any markdown fences the model might wrap around the JSON.
func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || end < start {
		return s
	}
	return s[start : end+1]
}
