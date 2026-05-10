package ai

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

var jsonBlock = regexp.MustCompile(`\{[\s\S]*\}`)

func ParseAndValidate(raw string) (VisionResponse, error) {
	clean := strings.TrimSpace(raw)
	if !strings.HasPrefix(clean, "{") {
		if m := jsonBlock.FindString(clean); m != "" {
			clean = m
		}
	}
	var out VisionResponse
	if err := json.Unmarshal([]byte(clean), &out); err != nil {
		return VisionResponse{}, err
	}
	if out.CropType == "" {
		return VisionResponse{}, errors.New("missing crop_type")
	}
	if len(out.Predictions) == 0 {
		return VisionResponse{}, errors.New("missing predictions")
	}
	for _, p := range out.Predictions {
		if p.Disease == "" {
			return VisionResponse{}, errors.New("missing disease")
		}
		if p.Confidence < 0 || p.Confidence > 1 {
			return VisionResponse{}, errors.New("invalid confidence")
		}
	}
	return out, nil
}
