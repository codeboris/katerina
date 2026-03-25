package ports

import "context"

type LLMResponse struct {
	Corrected   string `json:"corrected"`
	Translation string `json:"translation"`
	Answer      string `json:"answer"`
}

type LLMPort interface {
	Process(ctx context.Context, text string) (*LLMResponse, error)
}
