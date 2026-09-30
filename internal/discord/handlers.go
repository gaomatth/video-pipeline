package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// ignore messages sent by bots
	if m.Author.Bot {
		return
	}

	// checking attachments for file extensions
	for _, attachment := range m.Attachments {

		var filename string = attachment.Filename
		if strings.HasSuffix(filename, ".mp4") {
			// video attachment found! send job request to rabbitmq
		}
	}

}
