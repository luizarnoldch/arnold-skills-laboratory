package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	LabAPIURL        string
	OrchAPIURL       string
	WorkspaceDefault string
	HTTPTimeout      time.Duration
}

func LoadDotEnv() {
	_ = godotenv.Load()
}

func Load() (Config, error) {
	LoadDotEnv()

	timeout, err := time.ParseDuration(getenv("HTTP_TIMEOUT", "30s"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid HTTP_TIMEOUT: %w", err)
	}

	return Config{
		LabAPIURL:        trimTrailingSlash(getenv("LAB_API_URL", "http://127.0.0.1:18180")),
		OrchAPIURL:       trimTrailingSlash(getenv("ORCH_API_URL", "http://127.0.0.1:18182")),
		WorkspaceDefault: getenv("WORKSPACE_DEFAULT", "."),
		HTTPTimeout:      timeout,
	}, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func trimTrailingSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}
