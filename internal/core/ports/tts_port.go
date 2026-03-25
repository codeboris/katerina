package ports

import "context"

type TTSPort interface {
	Synthesize(ctx context.Context, text, voice string) (string, error)
}
