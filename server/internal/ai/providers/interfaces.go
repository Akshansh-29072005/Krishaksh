package providers

import "context"

type VisionProvider interface {
	Name() string
	Infer(ctx context.Context, prompt, cropType string, imageBytes []byte, mimeType string) (string, error)
}
