package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
}

// NewBot acts as a constructor, configuring the isolated session state
func NewBot(token string) (*Bot, error) {
	// Discord tokens must be prefixed with "Bot " when authenticating via WebSockets
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating discord session: %w", err)
	}

	bot := &Bot{session: dg}

	// Registering the message handler callback from handlers.go
	bot.session.AddHandler(bot.handleMessageCreate)

	// Instruct the session to only stream message events (limits network overhead)
	bot.session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	return bot, nil
}

// Connect opens the persistent outbound WebSocket pipeline to Discord
func (b *Bot) Connect() error {
	log.Println("Opening persistent WebSocket connection to Discord...")
	return b.session.Open()
}

// Close gracefully disconnects the WebSocket streamß
func (b *Bot) Close() error {
	log.Println("Closing Discord WebSocket session...")
	return b.session.Close()
}
