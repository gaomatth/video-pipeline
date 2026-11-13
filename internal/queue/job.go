package queue

import (
	"time"

	"github.com/google/uuid"
)

type VideoJob struct {
	JobId          uuid.UUID `json:"job_id"`
	ChannelId      string    `json:"channel_id"`
	MessageId      string    `json:"message_id"`
	UserId         string    `json:"user_id"`
	VideoURL       string    `json:"video_url"`
	Filename       string    `json:"filename"`
	JobRequestedAt time.Time `json:"job_requested_at"`
}
