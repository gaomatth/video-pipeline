package config

import (
	"fmt"
	"os"
)

type Config struct {
	DiscordToken string
	RabbitMQURL  string
	S3Endpoint   string
	S3Bucket     string
}

// Load fetches configurations from the environment.
// This allows you to easily switch between your local Docker Compose 
// variables and real AWS/GCP production secrets later.
func Load() (*Config, error) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN is required")
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	return &Config{
		DiscordToken: token,
		RabbitMQURL:  rabbitURL,
		S3Endpoint:   getEnv("S3_ENDPOINT", "localhost:9000"),
		S3Bucket:     getEnv("S3_BUCKET", "processed-videos"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
