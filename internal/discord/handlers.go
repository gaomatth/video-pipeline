package discord

import (
	"context"
	"log"
	"strings"
	"time"
	"video-pipeline/internal/queue"

	"github.com/google/uuid"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// ignore messages sent by bots
	if m.Author.Bot {
		return
	}

	if len(m.Attachments) == 0 {
		return
	}

	// checking attachments for file extensions
	for _, attachment := range m.Attachments {

		var filename string = attachment.Filename
		if strings.HasSuffix(strings.ToLower(filename), ".mp4") {
			// video attachment found! create job to be published

			job := queue.VideoJob{
				JobId:          uuid.New(),
				ChannelId:      m.ChannelID,
				MessageId:      m.ID,
				UserId:         m.Author.ID,
				VideoURL:       attachment.URL,
				Filename:       filename,
				JobRequestedAt: time.Now().UTC(),
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := b.publisher.PublishVideoJob(ctx, job)

			if err != nil {
				log.Printf("[Gateway] Failed to publish job: %v", err)
				_, _ = s.ChannelMessageSend(m.ChannelID, "❌ Failed to queue video for processing.")
				continue
			}

			_, _ = s.ChannelMessageSend(m.ChannelID, "✅ Video successfully queued for transcoding!")
		}
	}

}
