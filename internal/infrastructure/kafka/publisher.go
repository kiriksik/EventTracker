package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type EventPublisher struct {
	writer *kafka.Writer
}

func NewEventPublisher(brokers []string, topic string) *EventPublisher {
	return &EventPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *EventPublisher) Publish(key string, value []byte) error {
	return p.writer.WriteMessages(
		context.TODO(),
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
