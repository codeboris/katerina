package ports

import "context"

type LLMResponse struct {
	Corrected         string `json:"corrected"`
	Translation       string `json:"translation"`
	Answer            string `json:"answer"`
	AnswerTranslation string `json:"answer_translation"`
}

type LLMPort interface {
	Process(ctx context.Context, text string) (*LLMResponse, error)
}
