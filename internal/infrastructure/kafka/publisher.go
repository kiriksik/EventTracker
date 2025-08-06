package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type EventPublisher struct {
	writer *kafka.Writer
}

func NewEventPublisher(brokers []string, topic string) *EventPublisher {
	return &EventPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
			Async:        false,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (p *EventPublisher) Publish(ctx context.Context, key string, value []byte) error {
	return p.writer.WriteMessages(
		ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: value,
		},
	)
}

func (p *EventPublisher) Close() {
	if err := p.writer.Close(); err != nil {
		log.Printf("failed to close kafka writer: %v", err)
	}
}
