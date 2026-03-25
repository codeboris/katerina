package ports

import "context"

type TTSPort interface {
	Synthesize(ctx context.Context, text string) (string, error)
}
