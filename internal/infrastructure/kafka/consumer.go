package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kiriksik/EventTracker/internal/domain/entity"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type EventConsumer struct {
	reader *kafka.Reader
	logger *zap.Logger
}

func NewEventConsumer(brokers []string, topic string, logger *zap.Logger) *EventConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        "event-consumer-group",
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	})

	return &EventConsumer{
		reader: r,
		logger: logger,
	}
}

func (c *EventConsumer) Consume(ctx context.Context) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			c.logger.Error("failed to read message", zap.Error(err))
			continue
		}

		var event entity.Event
		if err := json.Unmarshal(m.Value, &event); err != nil {
			c.logger.Error("failed to unmarshal event", zap.ByteString("raw", m.Value), zap.Error(err))
			continue
		}

		c.handleEvent(event)
	}
}

func (c *EventConsumer) handleEvent(event entity.Event) {
	c.logger.Info("Received event",
		zap.String("id", event.ID),
		zap.String("type", event.Type),
		zap.Time("timestamp", event.Timestamp),
	)

}

func (c *EventConsumer) Close() error {
	return c.reader.Close()
}
