package queue

import (
	"context"
	"fmt"
	"log"

	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

type Publisher struct {
	conn      *rmq.AmqpConnection
	publisher *rmq.Publisher // Updated from amqp.AmqpPublisher to the package's Publisher interface
}

func NewPublisher(ctx context.Context, amqpURL string) (*Publisher, error) {
	conn, err := rmq.Dial(ctx, amqpURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ via AMQP 1.0: %w", err)
	}

	_, err = conn.Management().DeclareQueue(ctx, &rmq.QuorumQueueSpecification{
		Name: "video_jobs",
	})
	if err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	pub, err := conn.NewPublisher(ctx, &rmq.QueueAddress{Queue: "video_jobs"}, nil)
	if err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("failed to create AMQP 1.0 publisher: %w", err)
	}

	log.Println("[Publisher] Connected via AMQP 1.0 and initialized publisher for 'video_jobs'.")

	return &Publisher{
		conn:      conn,
		publisher: pub,
	}, nil
}

func (pub *Publisher) PublishVideoJob(ctx context.Context, job VideoJob) error {
	
}
