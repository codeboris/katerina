package ports

import "context"

type STTPort interface {
	Transcribe(ctx context.Context, audioPath string) (string, error)
}
