package queue

import (
	"context"
	"encoding/json"
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
	// TODO: Marshal the VideoJob into a JSON to make it AMQP compatible, then submit to queue
	json, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal VideoJob struct into JSON payload")
	}

	queueMessage := rmq.NewMessage(json)

	response, err := pub.publisher.Publish(ctx, queueMessage)
	if err != nil {
		return fmt.Errorf("failed to publish AMQP job message to AMQP")
	}

	return nil
}
