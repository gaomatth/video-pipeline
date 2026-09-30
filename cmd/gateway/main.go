package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	// 2. Initialize Infrastructure Layer First
	// This ensures our RabbitMQ queues are established before events happen
	rabbitClient, err := queue.NewRabbitMQClient(rabbitURL, "video-jobs")
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ connection client: %v", err)
	}
	defer rabbitClient.Close() // Safe network socket cleanup on exit

	// 3. Initialize Gateway / Interface Layer
	// We inject the running rabbitClient straight into the bot configuration constructor
	bot, err := discord.NewBot(token, rabbitClient)
	if err != nil {
		log.Fatalf("Failed to initialize Discord bot controller: %v", err)
	}

	// 4. Connect Outbound to Discord WebSockets
	if err := bot.Connect(); err != nil {
		log.Fatalf("Failed to establish Discord gateway connection: %v", err)
	}
	defer bot.Close()

	log.Println("API Ingestion Gateway is online and listening. Press CTRL+C to exit.")

	// 5. Graceful Lifecycle Shutdown Block
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	<-stopSignal

	log.Println("Gracefully terminating API Ingestion Gateway node...")
}
