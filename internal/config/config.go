package config

import (
	"fmt"
	"os"
)

type Config struct {
	GitLabToken     string
	GitLabBaseURL   string
	GeminiAPIKey    string
	WebhookSecret   string
}

func Load() (*Config, error) {
	cfg := &Config{
		GitLabToken:     os.Getenv("GITLAB_TOKEN"),
		GitLabBaseURL:   os.Getenv("GITLAB_BASE_URL"),
		GeminiAPIKey:    os.Getenv("GEMINI_API_KEY"),
		WebhookSecret:   os.Getenv("WEBHOOK_SECRET"),
	}

	if cfg.GitLabToken == "" {
		return nil, fmt.Errorf("GITLAB_TOKEN environment variable is required")
	}
	
	if cfg.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	if cfg.GitLabBaseURL == "" {
		cfg.GitLabBaseURL = "https://gitlab.com"
	}

	return cfg, nil
}