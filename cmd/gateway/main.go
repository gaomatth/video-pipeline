package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"video-pipeline/internal/discord"
	"video-pipeline/internal/queue"
)

func main() {
	// 1. Structural Config Checks
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatalf("Critical Error: DISCORD_TOKEN environment variable is not set")
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	log.Println("Initializing distributed pipeline components...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 2. Initialize Infrastructure Layer First
	// This ensures our RabbitMQ queues are established before events happen
	rabbitClient, err := queue.NewRabbitMQClient(ctx, rabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ connection client: %v", err)
	}
	defer func(rabbitClient *queue.Publisher, ctx context.Context) {
		err := rabbitClient.Close(ctx)
		if err != nil {

		}
	}(rabbitClient, ctx) // Safe network socket cleanup on exit

	// 3. Initialize Gateway / Interface Layer
	// We inject the running rabbitClient straight into the bot configuration constructor
	bot, err := discord.NewBot(token)
	if err != nil {
		log.Fatalf("Failed to initialize Discord bot controller: %v", err)
	}

	// 4. Connect Outbound to Discord WebSockets
	if err := bot.Connect(); err != nil {
		log.Fatalf("Failed to establish Discord gateway connection: %v", err)
	}
	defer func(bot *discord.Bot) {
		err := bot.Close()
		if err != nil {
			log.Fatalf("Failed to close Discord gateway connection: %v", err)
		}
	}(bot)

	log.Println("API Ingestion Gateway is online and listening. Press CTRL+C to exit.")

	// 5. Graceful Lifecycle Shutdown Block
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	<-stopSignal

	log.Println("Gracefully terminating API Ingestion Gateway node...")
}
