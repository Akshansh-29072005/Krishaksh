package ai

import (
	"fmt"
	"os"
	"path/filepath"
)

type PromptManager struct {
	baseDir string
}

func NewPromptManager(baseDir string) *PromptManager { return &PromptManager{baseDir: baseDir} }

func (p *PromptManager) Get(version string) (string, error) {
	if version == "" {
		version = "vision_v1"
	}
	file := filepath.Join(p.baseDir, version+".txt")
	b, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("read prompt failed: %w", err)
	}
	return string(b), nil
}
