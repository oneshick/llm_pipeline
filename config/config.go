package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GroqAPIKey string
	Model      string
	MaxTokens  int
	CSVPath    string
	OutPath    string
	DelayMs    int
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		GroqAPIKey: getEnv("GROQ_API_KEY", ""),
		Model:      getEnv("GROQ_MODEL", "llama-3.1-8b-instant"),
		MaxTokens:  getEnvInt("GROQ_MAX_TOKENS", 512),
		CSVPath:    getEnv("CSV_PATH", "products.csv"),
		OutPath:    getEnv("OUT_PATH", "results.json"),
		DelayMs:    getEnvInt("DELAY_MS", 300),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}

	return fallback
}
